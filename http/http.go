package http

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/togglr-io/togglr"
	"go.uber.org/zap"
)

// Services define all of the injectable service interfaces used by the HTTP handlers
type Services struct {
	ToggleService togglr.ToggleService
}

// A Config captures all of the information necessary to setup an HTTP server
type Config struct {
	Host     string
	Port     uint
	Services Services
	Logger   *zap.Logger
}

// BuildRoutes creates a Router and binds HTTP handlers to the routes. Exported mostly for testing purposes, should
// call Listen with a Config for real use-cases
func BuildRoutes(cfg Config) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	// r.Use(Telemetry(cfg.Logger))
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000", "http://localhost:9001"},
	}))

	r.Post("/toggle", HandleTogglePost(cfg.Logger, cfg.Services.ToggleService))
	r.Get("/toggle", HandleToggleGet(cfg.Logger, cfg.Services.ToggleService))
	r.Get("/toggle/{id}", HandleToggleGetID(cfg.Logger, cfg.Services.ToggleService))
	r.Delete("/toggle/{id}", HandleToggleDelete(cfg.Logger, cfg.Services.ToggleService))

	return r
}

func Listen(cfg Config) error {
	return http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), BuildRoutes(cfg))
}
