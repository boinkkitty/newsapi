package logger

import (
	"context"
	"log/slog"
	"net/http"
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

// Before Middleware.
func AddLoggerMid(logger *slog.Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Inject logger into request context
		loggerCtx := ContextWithLogger(r.Context(), logger)

		// Clone request with new context containing logger
		r = r.Clone(loggerCtx)

		// Call next handler with updated request
		next.ServeHTTP(w, r)
	}
}

func LoggerMid(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := FromContext(r.Context())
		logger.Info("Request", "path", r.URL.String())
		next.ServeHTTP(w, r)
	}
}
