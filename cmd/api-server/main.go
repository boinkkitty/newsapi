package main

import (
	"github.com/boinkkitty/newsapi/internal/news"
	"github.com/boinkkitty/newsapi/internal/postgres"
	"log/slog"
	"net/http"
	"os"

	"github.com/boinkkitty/newsapi/internal/logger"
	"github.com/boinkkitty/newsapi/internal/router"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	db, err := postgres.NewDB(&postgres.Config{})
	if err != nil {
		log.Error("db error", "err", err)
		os.Exit(1)
	}
	newsStore := news.NewStore(db)

	r := router.NewRouter(newsStore)
	wrappedRouter := logger.AddLoggerMid(log, logger.LoggerMid(r))

	log.Info("server starting on port 5002")

	if err := http.ListenAndServe(":5002", wrappedRouter); err != nil {
		log.Error("Failed to start server", "error", err)
	}
}
