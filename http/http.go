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
	ToggleService   togglr.ToggleService
	MetadataService togglr.MetadataService
	AccountService  togglr.AccountService
	UserService     togglr.UserService
	Resolver        togglr.Resolver
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

	r.Post("/toggle", HandleTogglePOST(cfg.Logger, cfg.Services.ToggleService))
	r.Get("/toggle", HandleToggleGET(cfg.Logger, cfg.Services.ToggleService))
	r.Get("/toggle/{id}", HandleToggleIdGET(cfg.Logger, cfg.Services.ToggleService))
	r.Delete("/toggle/{id}", HandleToggleDELETE(cfg.Logger, cfg.Services.ToggleService))

	r.Get("/metadata/{accountID}", HandleMetadataGET(cfg.Logger, cfg.Services.MetadataService))
	r.Post("/resolve/{accountID}", HandleResolvePOST(cfg.Logger, cfg.Services.Resolver))

	r.Post("/account", HandleAccountPOST(cfg.Logger, cfg.Services.AccountService))
	r.Get("/account", HandleAccountGET(cfg.Logger, cfg.Services.AccountService))
	r.Get("/account/{id}", HandleAccountIdGET(cfg.Logger, cfg.Services.AccountService))
	r.Get("/account/{id}/user", HandleAccountUsersGET(cfg.Logger, cfg.Services.UserService))
	r.Post("/account/{id}/user", HandleAccountUsersPOST(cfg.Logger, cfg.Services.AccountService))

	r.Post("/user", HandleUserPOST(cfg.Logger, cfg.Services.UserService))
	// a GET on /user returns the currently logged in user
	r.Get("/user", HandleUserGET(cfg.Logger, cfg.Services.UserService))
	r.Get("/user/{id}", HandleUserIdGET(cfg.Logger, cfg.Services.UserService))
	r.Delete("/user/{id}", HandleUserDELETE(cfg.Logger, cfg.Services.UserService))

	// oauth stuff
	r.Get("/oauth/github", HandleGithubRedirect(cfg.Logger, cfg.Services.UserService))

	return r
}

func Listen(cfg Config) error {
	return http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), BuildRoutes(cfg))
}
