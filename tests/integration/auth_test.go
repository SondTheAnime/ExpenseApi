package integration

import (
	"bytes"
	"encoding/json"
	"expenseapi/internal/auth"
	"expenseapi/internal/handler"
	"expenseapi/internal/model"
	"expenseapi/internal/repository"
	"expenseapi/internal/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

type testServer struct {
	db          *pgxpool.Pool
	authHandler *handler.AuthHandler
}

func setupTestServer(t *testing.T, db *pgxpool.Pool) *testServer {
	t.Helper()

	// Inicializa os serviços
	jwtService := auth.NewJWTService("test_secret_key", 24*60*60, 7*24*60*60)
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, jwtService)
	authHandler := handler.NewAuthHandler(authService)

	return &testServer{
		db:          db,
		authHandler: authHandler,
	}
}

func TestRegister(t *testing.T) {
	// Limpa o banco antes dos testes
	if err := cleanDatabase(); err != nil {
		t.Fatalf("erro ao limpar banco de dados: %v", err)
	}

	tests := []struct {
		name          string
		input         model.CreateUserInput
		expectedCode  int
		expectedError string
	}{
		{
			name: "registro_valido",
			input: model.CreateUserInput{
				Email:    "test@example.com",
				Password: "password123",
			},
			expectedCode: http.StatusCreated,
		},
		{
			name: "email_invalido",
			input: model.CreateUserInput{
				Email:    "invalid-email",
				Password: "password123",
			},
			expectedCode:  http.StatusBadRequest,
			expectedError: "email inválido",
		},
		{
			name: "senha_curta",
			input: model.CreateUserInput{
				Email:    "test@example.com",
				Password: "123",
			},
			expectedCode:  http.StatusBadRequest,
			expectedError: "senha deve ter no mínimo 6 caracteres",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Configura o servidor de teste
			srv := setupTestServer(t, testDB)

			// Cria a requisição
			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Executa a requisição
			w := httptest.NewRecorder()
			srv.authHandler.Register(w, req)

			// Verifica o código de status
			if w.Code != tt.expectedCode {
				t.Errorf("código de status esperado %d, obtido %d", tt.expectedCode, w.Code)
			}

			// Se espera erro, verifica a mensagem
			if tt.expectedError != "" {
				var response map[string]interface{}
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatalf("erro ao decodificar resposta: %v", err)
				}

				if msg, ok := response["message"].(string); !ok || msg != tt.expectedError {
					t.Errorf("mensagem de erro esperada %q, obtida %q", tt.expectedError, msg)
				}
			}
		})
	}
}

func TestLogin(t *testing.T) {
	// Limpa o banco antes dos testes
	if err := cleanDatabase(); err != nil {
		t.Fatalf("erro ao limpar banco de dados: %v", err)
	}

	// Configura o servidor de teste
	srv := setupTestServer(t, testDB)

	// Cria um usuário para teste
	input := model.CreateUserInput{
		Email:    "test@example.com",
		Password: "password123",
	}

	// Registra o usuário
	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	srv.authHandler.Register(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("erro ao criar usuário de teste: código %d", w.Code)
	}

	tests := []struct {
		name          string
		input         model.LoginInput
		expectedCode  int
		expectedError string
	}{
		{
			name: "login_valido",
			input: model.LoginInput{
				Email:    "test@example.com",
				Password: "password123",
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "senha_incorreta",
			input: model.LoginInput{
				Email:    "test@example.com",
				Password: "wrong_password",
			},
			expectedCode:  http.StatusUnauthorized,
			expectedError: "credenciais inválidas",
		},
		{
			name: "email_nao_encontrado",
			input: model.LoginInput{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			expectedCode:  http.StatusUnauthorized,
			expectedError: "credenciais inválidas",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Cria a requisição
			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Executa a requisição
			w := httptest.NewRecorder()
			srv.authHandler.Login(w, req)

			// Verifica o código de status
			if w.Code != tt.expectedCode {
				t.Errorf("código de status esperado %d, obtido %d", tt.expectedCode, w.Code)
			}

			// Se espera erro, verifica a mensagem
			if tt.expectedError != "" {
				var response map[string]interface{}
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatalf("erro ao decodificar resposta: %v", err)
				}

				if msg, ok := response["message"].(string); !ok || msg != tt.expectedError {
					t.Errorf("mensagem de erro esperada %q, obtida %q", tt.expectedError, msg)
				}
			}

			// Se login válido, verifica se retornou os tokens
			if tt.expectedCode == http.StatusOK {
				var response model.LoginResponse
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatalf("erro ao decodificar resposta: %v", err)
				}

				if response.Token == "" {
					t.Error("token de acesso não retornado")
				}
				if response.RefreshToken == "" {
					t.Error("refresh token não retornado")
				}
			}
		})
	}
}
