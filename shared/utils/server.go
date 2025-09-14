package utils

import (
	"context"
	"fmt"
	"net/http"
)

type ErrorHandler func(w http.ResponseWriter, r *http.Request, err error)

type ServerOption func(*Server)

type Server struct {
	port         int
	router       *Router
	middlewares  []Middleware
	errorHandler ErrorHandler
	server       *http.Server
}

func NewServer(port int) *Server {
	return &Server{
		port:        port,
		router:      NewRouter(),
		middlewares: []Middleware{},
		// errorHandler: nil,
		// server: &http.Server{
		// 	Addr:    ":" + string(port),
		// 	Handler: nil, // Will be set later
		// },
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) Run() error {
	addr := fmt.Sprintf(":%d", s.port)
	s.server = &http.Server{
		Addr:    addr,
		Handler: s,
	}
	fmt.Printf("Server starting on port %d\n", s.port)
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) GetRouter() *Router {
	return s.router
}

func WithErrorHandler(handler ErrorHandler) ServerOption {
	return func(s *Server) {
		s.errorHandler = handler
	}
}

func DefaultErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}