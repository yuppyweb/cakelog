package adapter

import (
	"context"
	"log/slog"

	"github.com/yuppyweb/cakelog"
)

// Is the default key under which context arguments will be stored in slog entries.
const DefaultSlogArgsKey = "context"

// Is an adapter that allows using a slog.Logger as a cakelog.Logger.
type SlogLogger struct {
	// The underlying slog.Logger to which log messages will be forwarded.
	Logger *slog.Logger

	// The key under which the context arguments will be stored in slog entries.
	ArgsKey string
}

// Creates a new SlogLogger that wraps the provided slog.Logger.
func NewSlogLogger(logger *slog.Logger) *SlogLogger {
	return &SlogLogger{
		Logger:  logger,
		ArgsKey: DefaultSlogArgsKey,
	}
}

// Sends a debug message to the underlying slog.Logger with the provided context and arguments.
func (sl *SlogLogger) Debug(ctx context.Context, msg string, args ...any) {
	sl.Logger.DebugContext(ctx, msg, slog.Any(sl.ArgsKey, args))
}

// Sends an info message to the underlying slog.Logger with the provided context and arguments.
func (sl *SlogLogger) Info(ctx context.Context, msg string, args ...any) {
	sl.Logger.InfoContext(ctx, msg, slog.Any(sl.ArgsKey, args))
}

// Sends a warning message to the underlying slog.Logger with the provided context and arguments.
func (sl *SlogLogger) Warn(ctx context.Context, msg string, args ...any) {
	sl.Logger.WarnContext(ctx, msg, slog.Any(sl.ArgsKey, args))
}

// Sends an error message to the underlying slog.Logger with the provided context, error, and arguments.
func (sl *SlogLogger) Error(ctx context.Context, err error, args ...any) {
	sl.Logger.ErrorContext(ctx, err.Error(), slog.Any(sl.ArgsKey, args))
}

// Ensures that SlogLogger implements the cakelog.Logger interface.
var _ cakelog.Logger = (*SlogLogger)(nil)
