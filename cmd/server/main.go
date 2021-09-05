package main

import (
	"fmt"
	"log"

	"github.com/mattn/go-colorable"
	"github.com/togglr-io/togglr"
	"github.com/togglr-io/togglr/env"
	"github.com/togglr-io/togglr/http"
	"github.com/togglr-io/togglr/pg"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func buildLogger() *zap.Logger {
	logCfg := zap.NewDevelopmentEncoderConfig()
	logCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(logCfg),
		zapcore.AddSync(colorable.NewColorableStdout()),
		zapcore.DebugLevel,
	))
}

func run() error {
	log := buildLogger()
	defer log.Sync()

	// load configs
	host := env.GetString("TOGGLE_HOST", "localhost")
	port := env.GetUint("TOGGLE_PORT", 9001)

	// initialize postgres
	db, err := pg.NewClient(pg.ConfigFromEnv("TOGGLE"))
	if err != nil {
		return fmt.Errorf("failed to create database connection: %w", err)
	}

	// initialize services to be used
	services := http.Services{
		ToggleService:   togglr.NewToggleService(db, db, log),
		MetadataService: db,
		AccountService:  db,
		UserService:     db,
		Resolver:        togglr.NewResolver(db),
	}

	// build server
	cfg := http.Config{
		Host:     host,
		Port:     port,
		Logger:   log,
		Services: services,
	}

	log.Info("starting server", zap.String("host", host), zap.Uint("port", port))
	return http.Listen(cfg)
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
