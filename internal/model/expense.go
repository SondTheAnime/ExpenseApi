package model

import (
	"time"
)

// Category representa as categorias possíveis de despesas
type Category string

const (
	CategoryGroceries   Category = "MANTIMENTOS"
	CategoryLeisure     Category = "LAZER"
	CategoryElectronics Category = "ELETRONICA"
	CategoryUtilities   Category = "UTILITARIOS"
	CategoryClothing    Category = "ROUPAS"
	CategoryHealth      Category = "SAUDE"
	CategoryOthers      Category = "OUTROS"
)

// Expense representa uma despesa no sistema
type Expense struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	Category    Category  `json:"category"`
	Date        time.Time `json:"date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateExpenseInput representa os dados necessários para criar uma nova despesa
type CreateExpenseInput struct {
	Amount      float64  `json:"amount" validate:"required,gt=0"`
	Description string   `json:"description" validate:"required,min=3,max=255"`
	Category    Category `json:"category" validate:"required,oneof=MANTIMENTOS LAZER ELETRONICA UTILITARIOS ROUPAS SAUDE OUTROS"`
	Date        string   `json:"date" validate:"required,datetime=2006-01-02"`
}

// UpdateExpenseInput representa os dados que podem ser atualizados em uma despesa
type UpdateExpenseInput struct {
	Amount      *float64  `json:"amount,omitempty" validate:"omitempty,gt=0"`
	Description *string   `json:"description,omitempty" validate:"omitempty,min=3,max=255"`
	Category    *Category `json:"category,omitempty" validate:"omitempty,oneof=MANTIMENTOS LAZER ELETRONICA UTILITARIOS ROUPAS SAUDE OUTROS"`
	Date        *string   `json:"date,omitempty" validate:"omitempty,datetime=2006-01-02"`
}

// ExpenseFilter representa os filtros disponíveis para busca de despesas
type ExpenseFilter struct {
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	Category  *Category  `json:"category,omitempty"`
}
