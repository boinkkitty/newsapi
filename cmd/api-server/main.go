package main

import (
	"github.com/boinkkitty/newsapi/internal/router"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	logger.Info("server starting on port 5002")

	r := router.NewRouter()
	if err := http.ListenAndServe(":5002", r); err != nil {
		logger.Error("Failed to start server", "error", err)
	}
}
