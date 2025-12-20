package main

import (
	"bearury/rest-api/internal/config"
	"bearury/rest-api/internal/lib/logger/sl"
	"bearury/rest-api/internal/storage/sqlite"
	"fmt"
	"log/slog"
	_ "log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// TODO init config: cleanenv
	cfg := config.MustLoad()

	fmt.Println(cfg)

	// TODO init logger : slog
	log := setupLogger(cfg.Environment)

	log.Info("starting server", slog.String("environment", cfg.Environment))
	log.Debug("debug logging enabled")

	// TODO init storage: sglite

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	_ = storage
}

func setupLogger(env string) *slog.Logger {

	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
