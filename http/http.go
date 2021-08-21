package http

import (
	"fmt"
	"net/http"

	"github.com/eriktate/toggle"
	"github.com/eriktate/toggle/uid"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

// Services define all of the injectable service interfaces used by the HTTP handlers
type Services struct {
	ToggleService toggle.ToggleService
}

// A Config captures all of the information necessary to setup an HTTP server
type Config struct {
	Host     string
	Port     uint
	Services Services
	Logger   *zap.Logger
}

type idEnvelope struct {
	ID uid.UID `json:"id"`
}

func buildRoutes(cfg Config) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	// r.Use(Telemetry(cfg.Logger))
	r.Use(middleware.Recoverer)

	r.Post("/toggle", handleTogglePost(cfg.Logger, cfg.Services.ToggleService))
	r.Get("/toggle", handleToggleList(cfg.Logger, cfg.Services.ToggleService))
	r.Get("/toggle/{id}", handleToggleDetail(cfg.Logger, cfg.Services.ToggleService))
	r.Delete("/toggle/{id}", handleToggleDelete(cfg.Logger, cfg.Services.ToggleService))

	return r
}

func Listen(cfg Config) error {
	return http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), buildRoutes(cfg))
}
