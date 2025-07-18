package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/boinkkitty/newsapi/internal/logger"
	"github.com/boinkkitty/newsapi/internal/router"
	"github.com/boinkkitty/newsapi/internal/store"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	r := router.NewRouter(store.New())
	wrappedRouter := logger.AddLoggerMid(log, logger.LoggerMid(r))

	log.Info("server starting on port 5002")

	if err := http.ListenAndServe(":5002", wrappedRouter); err != nil {
		log.Error("Failed to start server", "error", err)
	}
}
