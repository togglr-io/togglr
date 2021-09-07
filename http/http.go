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
	Host       string
	Port       uint
	HMACSecret []byte
	Services   Services
	Logger     *zap.Logger
}

func makeBaseRouter(cfg Config) chi.Router {
	r := chi.NewRouter()
	return r
}

// BuildRoutes creates a Router and binds HTTP handlers to the routes. Exported mostly for testing purposes, should
// call Listen with a Config for real use-cases
func BuildRoutes(cfg Config) chi.Router {
	tokenFn := makeTokenFn(cfg.HMACSecret)
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(Telemetry(cfg.Logger))
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:9001"},
		AllowCredentials: true,
	}))

	// unauthenticated endpoints
	r.Get("/oauth/github", HandleGithubRedirect(cfg.Logger, cfg.Services.UserService, tokenFn))

	authd := chi.NewRouter()
	authd.Use(requireAuth(cfg.HMACSecret, cfg.Logger))

	authd.Post("/toggle", HandleTogglePOST(cfg.Logger, cfg.Services.ToggleService))
	authd.Get("/toggle", HandleToggleGET(cfg.Logger, cfg.Services.ToggleService))
	authd.Get("/toggle/{id}", HandleToggleIdGET(cfg.Logger, cfg.Services.ToggleService))
	authd.Delete("/toggle/{id}", HandleToggleDELETE(cfg.Logger, cfg.Services.ToggleService))

	authd.Post("/account", HandleAccountPOST(cfg.Logger, cfg.Services.AccountService))
	authd.Get("/account", HandleAccountGET(cfg.Logger, cfg.Services.AccountService))
	authd.Get("/account/{id}", HandleAccountIdGET(cfg.Logger, cfg.Services.AccountService))
	authd.Get("/account/{id}/user", HandleAccountUsersGET(cfg.Logger, cfg.Services.UserService))
	authd.Post("/account/{id}/user", HandleAccountUsersPOST(cfg.Logger, cfg.Services.AccountService))

	authd.Post("/user", HandleUserPOST(cfg.Logger, cfg.Services.UserService))
	authd.Get("/user", HandleUserGET(cfg.Logger, cfg.Services.UserService))
	authd.Get("/user/{id}", HandleUserIdGET(cfg.Logger, cfg.Services.UserService))
	authd.Delete("/user/{id}", HandleUserDELETE(cfg.Logger, cfg.Services.UserService))

	authd.Get("/metadata/{accountID}", HandleMetadataGET(cfg.Logger, cfg.Services.MetadataService))

	authd.Post("/resolve/{accountID}", HandleResolvePOST(cfg.Logger, cfg.Services.Resolver))

	r.Mount("/", authd)
	return r
}

func Listen(cfg Config) error {
	return http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), BuildRoutes(cfg))
}
