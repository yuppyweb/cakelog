package zapadapter

import (
	"context"
	"errors"

	"github.com/yuppyweb/cakelog"
	"github.com/yuppyweb/cakelog/adapter"
	"go.uber.org/zap"
)

var ErrNilZapLogger = errors.New("zap logger is nil")

// zapAdapter is an adapter that allows using a zap.Logger as a cakelog.Logger.
type zapAdapter struct {
	// The underlying zap.Logger to which log messages will be forwarded.
	log *zap.Logger
}

// New creates a new zapAdapter that implements the cakelog.Logger interface.
func New(logger *zap.Logger) (cakelog.Logger, error) {
	if logger == nil {
		return nil, ErrNilZapLogger
	}

	return &zapAdapter{log: logger}, nil
}

// Debug sends a debug message to the underlying zap.Logger with the provided arguments.
// Note: context.Context is not used because zap.Logger does not have built-in context propagation
// for standard logging methods (unlike slog). To use context with zap, configure a logger hook or middleware.
func (zl *zapAdapter) Debug(_ context.Context, msg string, args ...any) {
	zl.log.Debug(msg, zapFields(args)...)
}

// Info sends an info message to the underlying zap.Logger with the provided arguments.
// Note: context.Context is not used because zap.Logger does not have built-in context propagation
// for standard logging methods (unlike slog). To use context with zap, configure a logger hook or middleware.
func (zl *zapAdapter) Info(_ context.Context, msg string, args ...any) {
	zl.log.Info(msg, zapFields(args)...)
}

// Warn sends a warning message to the underlying zap.Logger with the provided arguments.
// Note: context.Context is not used because zap.Logger does not have built-in context propagation
// for standard logging methods (unlike slog). To use context with zap, configure a logger hook or middleware.
func (zl *zapAdapter) Warn(_ context.Context, msg string, args ...any) {
	zl.log.Warn(msg, zapFields(args)...)
}

// Error sends an error message to the underlying zap.Logger with the provided error and arguments.
// Note: context.Context is not used because zap.Logger does not have built-in context propagation
// for standard logging methods (unlike slog). To use context with zap, configure a logger hook or middleware.
func (zl *zapAdapter) Error(_ context.Context, err error, args ...any) {
	msg := ""
	if err != nil {
		msg = err.Error()
	}

	zl.log.Error(msg, zapFields(args)...)
}

func zapFields(args []any) []zap.Field {
	fields := adapter.Fields(args)
	out := make([]zap.Field, len(fields))

	for i, field := range fields {
		out[i] = zap.Any(field.Key, field.Value)
	}

	return out
}

// Ensures that ZapLogger implements the cakelog.Logger interface.
var _ cakelog.Logger = (*zapAdapter)(nil)
