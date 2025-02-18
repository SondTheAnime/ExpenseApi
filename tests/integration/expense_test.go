package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"expenseapi/internal/auth"
	"expenseapi/internal/config"
	"expenseapi/internal/handler"
	"expenseapi/internal/middleware"
	"expenseapi/internal/model"
	"expenseapi/internal/repository"
	"expenseapi/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type expenseTestServer struct {
	db          *pgxpool.Pool
	jwtService  *auth.JWTService
	authToken   string
	userID      string
	httpHandler http.Handler
}

func setupExpenseTestServer(t *testing.T) *expenseTestServer {
	cfg := config.New()
	cfg.DB.Name = "expense_test_db" // Use o banco de dados de teste

	// Conecta ao banco de dados de teste
	dbpool, err := pgxpool.New(context.Background(), cfg.DB.DSN())
	require.NoError(t, err)

	// Limpa as tabelas antes dos testes
	_, err = dbpool.Exec(context.Background(), "TRUNCATE TABLE expenses CASCADE")
	require.NoError(t, err)
	_, err = dbpool.Exec(context.Background(), "TRUNCATE TABLE users CASCADE")
	require.NoError(t, err)

	// Configura os serviços
	jwtService := auth.NewJWTService(cfg.JWT.Secret, cfg.JWT.ExpiresIn, cfg.JWT.RefreshToken)
	userRepo := repository.NewUserRepository(dbpool)
	authService := service.NewAuthService(userRepo, jwtService)
	authHandler := handler.NewAuthHandler(authService)

	expenseRepo := repository.NewExpenseRepository(dbpool)
	expenseService := service.NewExpenseService(expenseRepo)
	expenseHandler := handler.NewExpenseHandler(expenseService)

	// Cria um usuário de teste
	email := "test@example.com"
	password := "password123"
	user, err := authService.Register(context.Background(), model.CreateUserInput{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)

	// Gera o token JWT para o usuário
	token, err := authService.Login(context.Background(), model.LoginInput{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)

	// Configura o router
	mux := http.NewServeMux()

	// Rotas de autenticação
	mux.HandleFunc("POST /api/v1/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)

	// Rotas de despesas (protegidas por autenticação)
	mux.HandleFunc("POST /api/v1/expenses", middleware.AuthMiddleware(jwtService, expenseHandler.Create))
	mux.HandleFunc("GET /api/v1/expenses", middleware.AuthMiddleware(jwtService, expenseHandler.List))
	mux.HandleFunc("GET /api/v1/expenses/{id}", middleware.AuthMiddleware(jwtService, expenseHandler.GetByID))
	mux.HandleFunc("PUT /api/v1/expenses/{id}", middleware.AuthMiddleware(jwtService, expenseHandler.Update))
	mux.HandleFunc("DELETE /api/v1/expenses/{id}", middleware.AuthMiddleware(jwtService, expenseHandler.Delete))

	return &expenseTestServer{
		db:          dbpool,
		jwtService:  jwtService,
		authToken:   token.Token,
		userID:      user.ID,
		httpHandler: mux,
	}
}

func TestExpenseEndpoints(t *testing.T) {
	server := setupExpenseTestServer(t)
	defer server.db.Close()

	t.Run("deve criar uma nova despesa", func(t *testing.T) {
		input := model.CreateExpenseInput{
			Amount:      150.50,
			Description: "Compras do mês",
			Category:    model.CategoryGroceries,
			Date:        time.Now().Format("2006-01-02"),
		}

		body, err := json.Marshal(input)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/expenses", bytes.NewReader(body))
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", server.authToken))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		server.httpHandler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response model.Expense
		err = json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)

		assert.NotEmpty(t, response.ID)
		assert.Equal(t, server.userID, response.UserID)
		assert.Equal(t, input.Amount, response.Amount)
		assert.Equal(t, input.Description, response.Description)
		assert.Equal(t, input.Category, response.Category)
	})

	t.Run("deve listar despesas do usuário", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/expenses", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", server.authToken))

		w := httptest.NewRecorder()
		server.httpHandler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []*model.Expense
		err := json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)

		assert.Len(t, response, 1)
	})

	t.Run("deve retornar erro ao tentar acessar sem autenticação", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/expenses", nil)
		w := httptest.NewRecorder()
		server.httpHandler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("deve atualizar uma despesa existente", func(t *testing.T) {
		// Primeiro, vamos obter a lista de despesas para pegar o ID da despesa criada
		req := httptest.NewRequest(http.MethodGet, "/api/v1/expenses", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", server.authToken))

		w := httptest.NewRecorder()
		server.httpHandler.ServeHTTP(w, req)

		var expenses []*model.Expense
		err := json.NewDecoder(w.Body).Decode(&expenses)
		require.NoError(t, err)
		require.NotEmpty(t, expenses)

		expenseID := expenses[0].ID
		newAmount := 175.50
		newDescription := "Compras do mês (atualizado)"

		input := model.UpdateExpenseInput{
			Amount:      &newAmount,
			Description: &newDescription,
		}

		body, err := json.Marshal(input)
		require.NoError(t, err)

		req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/expenses/%s", expenseID), bytes.NewReader(body))
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", server.authToken))
		req.Header.Set("Content-Type", "application/json")

		w = httptest.NewRecorder()
		server.httpHandler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var updated model.Expense
		err = json.NewDecoder(w.Body).Decode(&updated)
		require.NoError(t, err)

		assert.Equal(t, newAmount, updated.Amount)
		assert.Equal(t, newDescription, updated.Description)
	})

	t.Run("deve deletar uma despesa existente", func(t *testing.T) {
		// Primeiro, vamos obter a lista de despesas para pegar o ID da despesa
		req := httptest.NewRequest(http.MethodGet, "/api/v1/expenses", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", server.authToken))

		w := httptest.NewRecorder()
		server.httpHandler.ServeHTTP(w, req)

		var expenses []*model.Expense
		err := json.NewDecoder(w.Body).Decode(&expenses)
		require.NoError(t, err)
		require.NotEmpty(t, expenses)

		expenseID := expenses[0].ID

		req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/expenses/%s", expenseID), nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", server.authToken))

		w = httptest.NewRecorder()
		server.httpHandler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)

		// Verifica se a despesa foi realmente deletada
		req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/expenses/%s", expenseID), nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", server.authToken))

		w = httptest.NewRecorder()
		server.httpHandler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
