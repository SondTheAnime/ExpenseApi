package middleware

import (
	"context"
	"net/http"
	"strings"

	"expenseapi/internal/auth"
)

type contextKey string

const UserIDKey contextKey = "user_id"

// AuthMiddleware protege rotas que requerem autenticação
func AuthMiddleware(jwtService *auth.JWTService, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "token não fornecido", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "formato de token inválido", http.StatusUnauthorized)
			return
		}

		token := parts[1]
		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			http.Error(w, "token inválido", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// GetUserIDFromContext retorna o ID do usuário do contexto
func GetUserIDFromContext(ctx context.Context) string {
	if userID, ok := ctx.Value(UserIDKey).(string); ok {
		return userID
	}
	return ""
}
