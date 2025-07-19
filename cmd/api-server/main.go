package main

import (
	"context"
	"fmt"
	"github.com/boinkkitty/newsapi/internal/news"
	"github.com/boinkkitty/newsapi/internal/postgres"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/boinkkitty/newsapi/internal/logger"
	"github.com/boinkkitty/newsapi/internal/router"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	db, err := postgres.NewDB(&postgres.Config{
		Host:     os.Getenv("DATABASE_HOST"),
		DBName:   os.Getenv("DATABASE_NAME"),
		Password: os.Getenv("DATABASE_PASSWORD"),
		User:     os.Getenv("DATABASE_USER"),
		Port:     os.Getenv("DATABASE_PORT"),
		SSLMode:  "disable",
	})
	if err != nil {
		log.Error("db error", "err", err)
		os.Exit(1)
	}
	newsStore := news.NewStore(db)

	r := router.NewRouter(newsStore)
	wrappedRouter := logger.AddLoggerMid(log, logger.LoggerMid(r))

	log.Info("server starting on port 8080")

	server := &http.Server{
		Addr:              ":8080",
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           wrappedRouter,
	}

	errGrp, errGrpCtx := errgroup.WithContext(context.Background())
	errGrp.Go(func() error {
		if err := server.ListenAndServe(); err != nil {
			log.Error("failed to start server", "error", err)
			return fmt.Errorf("could not start server: %w", err)
		}
		return nil
	})

	errGrp.Go(func() error {
		sigch := make(chan os.Signal, 1)
		signal.Notify(sigch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		select {
		case sig := <-sigch:
			log.Info("received signal", "signal", sig)
		case <-errGrpCtx.Done():

		}

		ctxWithTimeout, cancelFn := context.WithTimeout(errGrpCtx, 5*time.Second)
		defer cancelFn()

		log.Info("initiate graceful shutdown")

		if err := server.Shutdown(ctxWithTimeout); err != nil {
			return fmt.Errorf("error graceful shutdown: %w", err)
		}

		return nil
	})

	if err := errGrp.Wait(); err != nil {
		log.Error("error running", "error", err)
	}
}
