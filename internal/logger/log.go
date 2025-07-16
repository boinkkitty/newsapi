package logger

import (
	"context"
	"log/slog"
	"os"
)

type CtxKey struct{}

func ContextWithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	if logger == nil {
		return ctx
	}

	// Return if already has logger
	if ctxLog, ok := ctx.Value(CtxKey{}).(*slog.Logger); ok && ctxLog == logger {
		return ctx
	}

	// Return new logger
	return context.WithValue(ctx, CtxKey{}, logger)
}

func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(CtxKey{}).(*slog.Logger); ok {
		return logger
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
}
