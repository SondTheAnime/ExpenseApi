package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"expenseapi/internal/auth"
	"expenseapi/internal/config"
	"expenseapi/internal/handler"
	"expenseapi/internal/middleware"
	"expenseapi/internal/repository"
	"expenseapi/internal/service"

	scalargo "github.com/bdpiprava/scalar-go"
	"github.com/bdpiprava/scalar-go/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Carrega as configurações
	cfg := config.New()

	// Conecta ao banco de dados
	dbpool, err := pgxpool.New(context.Background(), cfg.DB.DSN())
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer dbpool.Close()

	// Inicializa os serviços
	jwtService := auth.NewJWTService(cfg.JWT.Secret, cfg.JWT.ExpiresIn, cfg.JWT.RefreshToken)
	userRepo := repository.NewUserRepository(dbpool)
	authService := service.NewAuthService(userRepo, jwtService)
	authHandler := handler.NewAuthHandler(authService)

	// Inicializa os serviços de despesas
	expenseRepo := repository.NewExpenseRepository(dbpool)
	expenseService := service.NewExpenseService(expenseRepo)
	expenseHandler := handler.NewExpenseHandler(expenseService)

	// Configuração do router
	mux := http.NewServeMux()

	// Rota principal
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "API de Controle de Despesas - v1.0.0")
	})

	// Rotas de autenticação
	mux.HandleFunc("POST /api/v1/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)

	// Rotas de despesas (protegidas por autenticação)
	mux.HandleFunc("POST /api/v1/expenses", middleware.AuthMiddleware(jwtService, expenseHandler.Create))
	mux.HandleFunc("GET /api/v1/expenses", middleware.AuthMiddleware(jwtService, expenseHandler.List))
	mux.HandleFunc("GET /api/v1/expenses/{id}", middleware.AuthMiddleware(jwtService, expenseHandler.GetByID))
	mux.HandleFunc("PUT /api/v1/expenses/{id}", middleware.AuthMiddleware(jwtService, expenseHandler.Update))
	mux.HandleFunc("DELETE /api/v1/expenses/{id}", middleware.AuthMiddleware(jwtService, expenseHandler.Delete))

	// Rota para a documentação Scalar
	mux.HandleFunc("GET /docs", func(w http.ResponseWriter, r *http.Request) {
		docsDir := filepath.Join(".", "docs")
		content, err := scalargo.New(
			docsDir,
			scalargo.WithBaseFileName("api.yaml"),
			scalargo.WithSpecModifier(func(spec *model.Spec) *model.Spec {
				spec.Info.Title = "API de Controle de Despesas"
				return spec
			}),
		)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(content))
	})

	// Configuração do servidor
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: middleware.CORS(mux),
	}

	log.Printf("Iniciando servidor na porta %s", server.Addr)
	log.Printf("Documentação disponível em http://localhost%s/docs", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
