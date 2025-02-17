package unit

import (
	"expenseapi/internal/auth"
	"testing"
	"time"
)

func TestJWTService(t *testing.T) {
	// Configuração do serviço
	secretKey := "test_secret_key"
	expiresIn := 24 * time.Hour
	refreshExpiry := 7 * 24 * time.Hour
	jwtService := auth.NewJWTService(secretKey, expiresIn, refreshExpiry)

	t.Run("gerar_token", func(t *testing.T) {
		userID := "test-user-id"

		// Gera os tokens
		token, refreshToken, err := jwtService.GenerateToken(userID)
		if err != nil {
			t.Fatalf("erro ao gerar token: %v", err)
		}

		// Verifica se os tokens foram gerados
		if token == "" {
			t.Error("token de acesso não gerado")
		}
		if refreshToken == "" {
			t.Error("refresh token não gerado")
		}

		// Valida o token de acesso
		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			t.Fatalf("erro ao validar token: %v", err)
		}

		// Verifica as claims
		if claims.UserID != userID {
			t.Errorf("user_id esperado %q, obtido %q", userID, claims.UserID)
		}
	})

	t.Run("token_invalido", func(t *testing.T) {
		// Tenta validar um token inválido
		_, err := jwtService.ValidateToken("invalid.token.here")
		if err == nil {
			t.Error("esperado erro ao validar token inválido")
		}
	})

	t.Run("token_expirado", func(t *testing.T) {
		// Cria um serviço com token que expira em 1 segundo
		shortJWT := auth.NewJWTService(secretKey, 1*time.Second, refreshExpiry)

		// Gera o token
		token, _, err := shortJWT.GenerateToken("test-user-id")
		if err != nil {
			t.Fatalf("erro ao gerar token: %v", err)
		}

		// Espera o token expirar
		time.Sleep(2 * time.Second)

		// Tenta validar o token expirado
		_, err = shortJWT.ValidateToken(token)
		if err == nil {
			t.Error("esperado erro ao validar token expirado")
		}
	})
}
