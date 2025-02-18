package repository

import (
	"context"
	"errors"
	"time"

	"expenseapi/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ExpenseRepository é a interface que define os métodos do repositório de despesas
type ExpenseRepository interface {
	Create(ctx context.Context, expense *model.Expense) error
	GetByID(ctx context.Context, id string, userID string) (*model.Expense, error)
	List(ctx context.Context, userID string, filter *model.ExpenseFilter) ([]*model.Expense, error)
	Update(ctx context.Context, expense *model.Expense) error
	Delete(ctx context.Context, id string, userID string) error
}

// PostgresExpenseRepository gerencia o acesso aos dados de despesas no banco
type PostgresExpenseRepository struct {
	db *pgxpool.Pool
}

// NewExpenseRepository cria uma nova instância do repositório de despesas
func NewExpenseRepository(db *pgxpool.Pool) ExpenseRepository {
	return &PostgresExpenseRepository{db: db}
}

// Create insere uma nova despesa no banco de dados
func (r *PostgresExpenseRepository) Create(ctx context.Context, expense *model.Expense) error {
	query := `
		INSERT INTO expenses (id, user_id, amount, description, category, date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	expense.ID = uuid.New().String()
	expense.CreatedAt = time.Now()
	expense.UpdatedAt = expense.CreatedAt

	_, err := r.db.Exec(ctx, query,
		expense.ID,
		expense.UserID,
		expense.Amount,
		expense.Description,
		expense.Category,
		expense.Date,
		expense.CreatedAt,
		expense.UpdatedAt,
	)

	return err
}

// GetByID busca uma despesa pelo ID
func (r *PostgresExpenseRepository) GetByID(ctx context.Context, id string, userID string) (*model.Expense, error) {
	query := `
		SELECT id, user_id, amount, description, category, date, created_at, updated_at
		FROM expenses
		WHERE id = $1 AND user_id = $2
	`

	expense := &model.Expense{}
	err := r.db.QueryRow(ctx, query, id, userID).Scan(
		&expense.ID,
		&expense.UserID,
		&expense.Amount,
		&expense.Description,
		&expense.Category,
		&expense.Date,
		&expense.CreatedAt,
		&expense.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, errors.New("despesa não encontrada")
	}

	return expense, err
}

// List retorna todas as despesas de um usuário com filtros opcionais
func (r *PostgresExpenseRepository) List(ctx context.Context, userID string, filter *model.ExpenseFilter) ([]*model.Expense, error) {
	query := `
		SELECT id, user_id, amount, description, category, date, created_at, updated_at
		FROM expenses
		WHERE user_id = $1
	`
	args := []interface{}{userID}
	argCount := 2

	if filter != nil {
		if filter.StartDate != nil {
			query += ` AND date >= $` + string(rune('0'+argCount))
			args = append(args, filter.StartDate)
			argCount++
		}
		if filter.EndDate != nil {
			query += ` AND date <= $` + string(rune('0'+argCount))
			args = append(args, filter.EndDate)
			argCount++
		}
		if filter.Category != nil {
			query += ` AND category = $` + string(rune('0'+argCount))
			args = append(args, filter.Category)
			argCount++
		}
	}

	query += ` ORDER BY date DESC`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []*model.Expense
	for rows.Next() {
		expense := &model.Expense{}
		err := rows.Scan(
			&expense.ID,
			&expense.UserID,
			&expense.Amount,
			&expense.Description,
			&expense.Category,
			&expense.Date,
			&expense.CreatedAt,
			&expense.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, expense)
	}

	return expenses, nil
}

// Update atualiza uma despesa existente
func (r *PostgresExpenseRepository) Update(ctx context.Context, expense *model.Expense) error {
	query := `
		UPDATE expenses
		SET amount = $1, description = $2, category = $3, date = $4, updated_at = $5
		WHERE id = $6 AND user_id = $7
	`

	expense.UpdatedAt = time.Now()

	result, err := r.db.Exec(ctx, query,
		expense.Amount,
		expense.Description,
		expense.Category,
		expense.Date,
		expense.UpdatedAt,
		expense.ID,
		expense.UserID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("despesa não encontrada")
	}

	return nil
}

// Delete remove uma despesa do banco de dados
func (r *PostgresExpenseRepository) Delete(ctx context.Context, id string, userID string) error {
	query := `DELETE FROM expenses WHERE id = $1 AND user_id = $2`

	result, err := r.db.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("despesa não encontrada")
	}

	return nil
}
