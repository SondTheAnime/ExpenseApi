package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"expenseapi/internal/model"
	"expenseapi/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockExpenseRepository é um mock do repositório de despesas
type MockExpenseRepository struct {
	mock.Mock
}

func (m *MockExpenseRepository) Create(ctx context.Context, expense *model.Expense) error {
	args := m.Called(ctx, expense)
	return args.Error(0)
}

func (m *MockExpenseRepository) GetByID(ctx context.Context, id string, userID string) (*model.Expense, error) {
	args := m.Called(ctx, id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Expense), args.Error(1)
}

func (m *MockExpenseRepository) List(ctx context.Context, userID string, filter *model.ExpenseFilter) ([]*model.Expense, error) {
	args := m.Called(ctx, userID, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Expense), args.Error(1)
}

func (m *MockExpenseRepository) Update(ctx context.Context, expense *model.Expense) error {
	args := m.Called(ctx, expense)
	return args.Error(0)
}

func (m *MockExpenseRepository) Delete(ctx context.Context, id string, userID string) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func TestExpenseService_Create(t *testing.T) {
	mockRepo := new(MockExpenseRepository)
	service := service.NewExpenseService(mockRepo)
	ctx := context.Background()
	userID := "user123"

	t.Run("deve criar uma despesa com sucesso", func(t *testing.T) {
		input := &model.CreateExpenseInput{
			Amount:      100.50,
			Description: "Teste de despesa",
			Category:    model.CategoryGroceries,
			Date:        "2024-02-18",
		}

		mockRepo.On("Create", ctx, mock.AnythingOfType("*model.Expense")).Return(nil)

		expense, err := service.Create(ctx, userID, input)

		assert.NoError(t, err)
		assert.NotNil(t, expense)
		assert.Equal(t, userID, expense.UserID)
		assert.Equal(t, input.Amount, expense.Amount)
		assert.Equal(t, input.Description, expense.Description)
		assert.Equal(t, input.Category, expense.Category)

		mockRepo.AssertExpectations(t)
	})

	t.Run("deve retornar erro quando a data é inválida", func(t *testing.T) {
		input := &model.CreateExpenseInput{
			Amount:      100.50,
			Description: "Teste de despesa",
			Category:    model.CategoryGroceries,
			Date:        "data-invalida",
		}

		expense, err := service.Create(ctx, userID, input)

		assert.Error(t, err)
		assert.Nil(t, expense)
	})
}

func TestExpenseService_GetByID(t *testing.T) {
	mockRepo := new(MockExpenseRepository)
	service := service.NewExpenseService(mockRepo)
	ctx := context.Background()
	userID := "user123"
	expenseID := "expense123"

	t.Run("deve retornar uma despesa existente", func(t *testing.T) {
		expected := &model.Expense{
			ID:          expenseID,
			UserID:      userID,
			Amount:      100.50,
			Description: "Teste de despesa",
			Category:    model.CategoryGroceries,
			Date:        time.Now(),
		}

		mockRepo.On("GetByID", ctx, expenseID, userID).Return(expected, nil).Once()

		expense, err := service.GetByID(ctx, expenseID, userID)

		assert.NoError(t, err)
		assert.Equal(t, expected, expense)
		mockRepo.AssertExpectations(t)
	})

	t.Run("deve retornar erro quando a despesa não existe", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, expenseID, userID).Return(nil, errors.New("despesa não encontrada")).Once()

		expense, err := service.GetByID(ctx, expenseID, userID)

		assert.Error(t, err)
		assert.Nil(t, expense)
		mockRepo.AssertExpectations(t)
	})
}

func TestExpenseService_Update(t *testing.T) {
	mockRepo := new(MockExpenseRepository)
	service := service.NewExpenseService(mockRepo)
	ctx := context.Background()
	userID := "user123"
	expenseID := "expense123"

	t.Run("deve atualizar uma despesa com sucesso", func(t *testing.T) {
		existingExpense := &model.Expense{
			ID:          expenseID,
			UserID:      userID,
			Amount:      100.50,
			Description: "Despesa original",
			Category:    model.CategoryGroceries,
			Date:        time.Now(),
		}

		newAmount := 150.75
		newDescription := "Despesa atualizada"
		input := &model.UpdateExpenseInput{
			Amount:      &newAmount,
			Description: &newDescription,
		}

		mockRepo.On("GetByID", ctx, expenseID, userID).Return(existingExpense, nil)
		mockRepo.On("Update", ctx, mock.AnythingOfType("*model.Expense")).Return(nil)

		updated, err := service.Update(ctx, expenseID, userID, input)

		assert.NoError(t, err)
		assert.Equal(t, newAmount, updated.Amount)
		assert.Equal(t, newDescription, updated.Description)
		mockRepo.AssertExpectations(t)
	})
}

func TestExpenseService_Delete(t *testing.T) {
	mockRepo := new(MockExpenseRepository)
	service := service.NewExpenseService(mockRepo)
	ctx := context.Background()
	userID := "user123"
	expenseID := "expense123"

	t.Run("deve deletar uma despesa com sucesso", func(t *testing.T) {
		mockRepo.On("Delete", ctx, expenseID, userID).Return(nil).Once()

		err := service.Delete(ctx, expenseID, userID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("deve retornar erro quando a despesa não existe", func(t *testing.T) {
		mockRepo.On("Delete", ctx, expenseID, userID).Return(errors.New("despesa não encontrada")).Once()

		err := service.Delete(ctx, expenseID, userID)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestExpenseService_GetExpensesByPeriod(t *testing.T) {
	mockRepo := new(MockExpenseRepository)
	service := service.NewExpenseService(mockRepo)
	ctx := context.Background()
	userID := "user123"

	t.Run("deve retornar despesas do último mês", func(t *testing.T) {
		expected := []*model.Expense{
			{
				ID:          "expense1",
				UserID:      userID,
				Amount:      100.50,
				Description: "Despesa 1",
				Category:    model.CategoryGroceries,
				Date:        time.Now().AddDate(0, 0, -15),
			},
			{
				ID:          "expense2",
				UserID:      userID,
				Amount:      200.75,
				Description: "Despesa 2",
				Category:    model.CategoryGroceries,
				Date:        time.Now().AddDate(0, 0, -5),
			},
		}

		mockRepo.On("List", ctx, userID, mock.AnythingOfType("*model.ExpenseFilter")).Return(expected, nil)

		expenses, err := service.GetExpensesByPeriod(ctx, userID, "month")

		assert.NoError(t, err)
		assert.Equal(t, expected, expenses)
		mockRepo.AssertExpectations(t)
	})
}
