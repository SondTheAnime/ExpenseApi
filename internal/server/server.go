package server

import (
	"expenseapi/internal/handler"
	"expenseapi/internal/middleware"
	"net/http"
)

type Server struct {
	router *http.ServeMux
}

func NewServer(authHandler *handler.AuthHandler) *Server {
	server := &Server{
		router: http.NewServeMux(),
	}

	// Configuração das rotas
	server.router.HandleFunc("POST /api/v1/auth/register", authHandler.Register)
	server.router.HandleFunc("POST /api/v1/auth/login", authHandler.Login)

	return server
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Aplicar middleware CORS em todas as requisições
	handler := middleware.CORS(s.router)
	handler.ServeHTTP(w, r)
}

func (s *Server) Start(addr string) error {
	return http.ListenAndServe(addr, s)
}
