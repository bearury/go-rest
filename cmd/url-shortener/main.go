package main

import (
	"bearury/rest-api/internal/config"
	"bearury/rest-api/internal/http-server/handlers/redirect"
	"bearury/rest-api/internal/http-server/handlers/url/remove"
	"bearury/rest-api/internal/http-server/handlers/url/save"
	"bearury/rest-api/internal/http-server/middleware/logger"
	"bearury/rest-api/internal/lib/logger/handlers/slogpretty"
	"bearury/rest-api/internal/lib/logger/sl"
	"bearury/rest-api/internal/storage/sqlite"
	"fmt"
	"log/slog"
	_ "log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// init config: cleanenv
	cfg := config.MustLoad()

	fmt.Println(cfg)

	// init logger : slog
	log := setupLogger(cfg.Environment)

	log.Info("starting server", slog.String("environment", cfg.Environment))
	log.Debug("debug logging enabled")

	// init storage: sglite

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()
	// middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	//Для авторизации

	router.Post("/url", save.New(log, storage))
	router.Get("/{alias}", redirect.New(log, storage))
	router.Delete("/{alias}", remove.New(log, storage))

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", sl.Err(err))
	}

}

func setupLogger(env string) *slog.Logger {

	var log *slog.Logger

	switch env {
	case envLocal:
		//log = slog.New(
		//	slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		//)

		log = setupPrettySlog()
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

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}
	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
