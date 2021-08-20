package main

import (
	"log"
	"os"
	"strconv"

	"github.com/eriktate/toggle/http"
	"github.com/eriktate/toggle/mock"
	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func getEnvString(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return def
}

func getEnvUint(key string, def uint) uint {
	if val, err := strconv.ParseUint(os.Getenv(key), 10, 32); err == nil {
		return uint(val)
	}

	return def
}

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
	host := getEnvString("TOGGLE_HOST", "localhost")
	port := getEnvUint("TOGGLE_PORT", 9001)

	// initialize services to be used
	services := http.Services{
		ToggleService: mock.NewToggleService(nil),
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
