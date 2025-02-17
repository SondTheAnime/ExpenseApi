package integration

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"expenseapi/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	testConfig *config.Config
	testDB     *pgxpool.Pool
)

func TestMain(m *testing.M) {
	// Configuração para testes
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "expense_user")
	os.Setenv("DB_PASSWORD", "expense_password")
	os.Setenv("DB_NAME", "expense_test_db")
	os.Setenv("DB_SSL_MODE", "disable")
	os.Setenv("JWT_SECRET", "test_secret_key")

	// Carrega configurações
	testConfig = config.New()

	// Conecta ao banco de dados
	var err error
	testDB, err = pgxpool.New(context.Background(), testConfig.DB.DSN())
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados de teste: %v", err)
	}
	defer testDB.Close()

	// Limpa o banco de dados antes dos testes
	if err := cleanDatabase(); err != nil {
		log.Fatalf("Erro ao limpar banco de dados: %v", err)
	}

	// Executa os testes
	code := m.Run()

	// Limpa o banco de dados após os testes
	if err := cleanDatabase(); err != nil {
		log.Printf("Erro ao limpar banco de dados após os testes: %v", err)
	}

	os.Exit(code)
}

func cleanDatabase() error {
	queries := []string{
		"TRUNCATE users CASCADE",
		"TRUNCATE expenses CASCADE",
	}

	for _, query := range queries {
		if _, err := testDB.Exec(context.Background(), query); err != nil {
			return fmt.Errorf("erro ao executar query %s: %v", query, err)
		}
	}

	return nil
}
