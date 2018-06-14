// Package logcontext provides utilities to work with logger in the context.
package logcontext

import (
	"context"

	"go.uber.org/zap"
)

// logKey is the key, under which *zap.Logger
// is contained in the context.
type logKey struct{}

// WithLogger adds a logger to ctx.
func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, logKey{}, logger)
}

// Logger returns the *zap.Logger stored in ctx or zap.L() if no logger is stored.
func Logger(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(logKey{}).(*zap.Logger)
	if !ok || logger == nil {
		return zap.L()
	}
	return logger
}
