package logger_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/boinkkitty/newsapi/internal/logger"
)

func Test_ContextWithLogger(t *testing.T) {
	testCases := []struct {
		name         string
		ctx          context.Context
		logger       *slog.Logger
		loggerExists bool
	}{
		{
			name: "return context without logger",
			ctx:  context.Background(),
		},
		{
			name:         "return context as it is",
			ctx:          context.WithValue(context.Background(), logger.CtxKey{}, slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))),
			loggerExists: true,
		},
		{
			name:         "inject logger",
			ctx:          context.Background(),
			logger:       slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})),
			loggerExists: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := logger.ContextWithLogger(tc.ctx, tc.logger)

			_, ok := ctx.Value(logger.CtxKey{}).(*slog.Logger)
			if tc.loggerExists != ok {
				t.Errorf("got %v, want %v", tc.loggerExists, ok)
			}
		})
	}
}

func Test_FromContext(t *testing.T) {
	testCases := []struct {
		name           string
		ctx            context.Context
		expectedLogger bool
	}{
		{
			name:           "logger exists",
			ctx:            context.WithValue(context.Background(), logger.CtxKey{}, slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))),
			expectedLogger: true,
		},
		{
			name:           "return new logger",
			ctx:            context.Background(),
			expectedLogger: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			logger := logger.FromContext(tc.ctx)

			if tc.expectedLogger && logger == nil {
				t.Errorf("got %v, want %v", logger, tc.expectedLogger)
			}
		})
	}
}
