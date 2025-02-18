package middleware

import "net/http"

// CORS é um middleware que adiciona os headers necessários para permitir CORS
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Permitir origens específicas
		allowedOrigins := []string{
			"http://localhost:8080",
			"http://localhost:8081",
			"http://localhost:3000",
		}

		origin := r.Header.Get("Origin")
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		// Métodos HTTP permitidos
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")

		// Headers permitidos
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token")

		// Expor headers
		w.Header().Set("Access-Control-Expose-Headers", "Authorization")

		// Permitir credenciais
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Tempo de cache para preflight
		w.Header().Set("Access-Control-Max-Age", "86400")

		// Tratamento especial para requisições OPTIONS (preflight)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
