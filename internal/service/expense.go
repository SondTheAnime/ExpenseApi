package service

import (
	"context"
	"time"

	"expenseapi/internal/model"
	"expenseapi/internal/repository"
)

// ExpenseService gerencia a lógica de negócios relacionada a despesas
type ExpenseService struct {
	repo *repository.ExpenseRepository
}

// NewExpenseService cria uma nova instância do serviço de despesas
func NewExpenseService(repo *repository.ExpenseRepository) *ExpenseService {
	return &ExpenseService{repo: repo}
}

// Create cria uma nova despesa
func (s *ExpenseService) Create(ctx context.Context, userID string, input *model.CreateExpenseInput) (*model.Expense, error) {
	date, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		return nil, err
	}

	expense := &model.Expense{
		UserID:      userID,
		Amount:      input.Amount,
		Description: input.Description,
		Category:    input.Category,
		Date:        date,
	}

	err = s.repo.Create(ctx, expense)
	if err != nil {
		return nil, err
	}

	return expense, nil
}

// GetByID busca uma despesa pelo ID
func (s *ExpenseService) GetByID(ctx context.Context, id string, userID string) (*model.Expense, error) {
	return s.repo.GetByID(ctx, id, userID)
}

// List retorna todas as despesas de um usuário com filtros opcionais
func (s *ExpenseService) List(ctx context.Context, userID string, filter *model.ExpenseFilter) ([]*model.Expense, error) {
	return s.repo.List(ctx, userID, filter)
}

// Update atualiza uma despesa existente
func (s *ExpenseService) Update(ctx context.Context, id string, userID string, input *model.UpdateExpenseInput) (*model.Expense, error) {
	expense, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	if input.Amount != nil {
		expense.Amount = *input.Amount
	}
	if input.Description != nil {
		expense.Description = *input.Description
	}
	if input.Category != nil {
		expense.Category = *input.Category
	}
	if input.Date != nil {
		date, err := time.Parse("2006-01-02", *input.Date)
		if err != nil {
			return nil, err
		}
		expense.Date = date
	}

	err = s.repo.Update(ctx, expense)
	if err != nil {
		return nil, err
	}

	return expense, nil
}

// Delete remove uma despesa
func (s *ExpenseService) Delete(ctx context.Context, id string, userID string) error {
	return s.repo.Delete(ctx, id, userID)
}

// GetExpensesByPeriod retorna despesas filtradas por período
func (s *ExpenseService) GetExpensesByPeriod(ctx context.Context, userID string, period string) ([]*model.Expense, error) {
	now := time.Now()
	var startDate, endDate time.Time

	switch period {
	case "week":
		startDate = now.AddDate(0, 0, -7)
	case "month":
		startDate = now.AddDate(0, -1, 0)
	case "quarter":
		startDate = now.AddDate(0, -3, 0)
	default:
		return nil, nil
	}

	endDate = now

	filter := &model.ExpenseFilter{
		StartDate: &startDate,
		EndDate:   &endDate,
	}

	return s.repo.List(ctx, userID, filter)
}
