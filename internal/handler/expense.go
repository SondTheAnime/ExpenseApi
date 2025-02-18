package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"expenseapi/internal/middleware"
	"expenseapi/internal/model"
	"expenseapi/internal/service"
)

// ExpenseHandler gerencia as requisições HTTP relacionadas a despesas
type ExpenseHandler struct {
	service *service.ExpenseService
}

// NewExpenseHandler cria uma nova instância do handler de despesas
func NewExpenseHandler(service *service.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{service: service}
}

// Create cria uma nova despesa
func (h *ExpenseHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "usuário não autenticado", http.StatusUnauthorized)
		return
	}

	var input model.CreateExpenseInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "dados inválidos", http.StatusBadRequest)
		return
	}

	expense, err := h.service.Create(r.Context(), userID, &input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(expense)
}

// GetByID retorna uma despesa específica
func (h *ExpenseHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "usuário não autenticado", http.StatusUnauthorized)
		return
	}

	expenseID := r.PathValue("id")
	if expenseID == "" {
		http.Error(w, "ID da despesa não fornecido", http.StatusBadRequest)
		return
	}

	expense, err := h.service.GetByID(r.Context(), expenseID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expense)
}

// List retorna todas as despesas do usuário
func (h *ExpenseHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "usuário não autenticado", http.StatusUnauthorized)
		return
	}

	period := r.URL.Query().Get("period")
	var expenses []*model.Expense
	var err error

	if period != "" {
		expenses, err = h.service.GetExpensesByPeriod(r.Context(), userID, period)
	} else {
		startDateStr := r.URL.Query().Get("start_date")
		endDateStr := r.URL.Query().Get("end_date")
		categoryStr := r.URL.Query().Get("category")

		filter := &model.ExpenseFilter{}

		if startDateStr != "" {
			startDate, err := time.Parse("2006-01-02", startDateStr)
			if err == nil {
				filter.StartDate = &startDate
			}
		}

		if endDateStr != "" {
			endDate, err := time.Parse("2006-01-02", endDateStr)
			if err == nil {
				filter.EndDate = &endDate
			}
		}

		if categoryStr != "" {
			category := model.Category(categoryStr)
			filter.Category = &category
		}

		expenses, err = h.service.List(r.Context(), userID, filter)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenses)
}

// Update atualiza uma despesa existente
func (h *ExpenseHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "usuário não autenticado", http.StatusUnauthorized)
		return
	}

	expenseID := r.PathValue("id")
	if expenseID == "" {
		http.Error(w, "ID da despesa não fornecido", http.StatusBadRequest)
		return
	}

	var input model.UpdateExpenseInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "dados inválidos", http.StatusBadRequest)
		return
	}

	expense, err := h.service.Update(r.Context(), expenseID, userID, &input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expense)
}

// Delete remove uma despesa
func (h *ExpenseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, "usuário não autenticado", http.StatusUnauthorized)
		return
	}

	expenseID := r.PathValue("id")
	if expenseID == "" {
		http.Error(w, "ID da despesa não fornecido", http.StatusBadRequest)
		return
	}

	err := h.service.Delete(r.Context(), expenseID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
