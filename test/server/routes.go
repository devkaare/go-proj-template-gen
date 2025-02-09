package server

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/devkaare/foobar/handler"
	"github.com/devkaare/foobar/views"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/", templ.Handler(views.Hello()).ServeHTTP)

	r.Route("/api", s.RegisterNewRoutes)

	return r
}

func (s *Server) RegisterNewRoutes(r chi.Router) {
	handler := &handler.New{}

	r.Get("/hello", handler.Greet)

}

