package adapter

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/yuppyweb/cakelog"
)

var (
	ErrNilSlogLogger = errors.New("is nil slog.Logger")
	ErrNilSlogOption = errors.New("is nil slog option")
)

// SlogLogger is an adapter that allows using a slog.Logger as a cakelog.Logger.
type SlogLogger struct {
	// The underlying slog.Logger to which log messages will be forwarded.
	log *slog.Logger

	// The options for configuring the SlogLogger behavior.
	opt *Options
}

// NewSlogLogger creates a new SlogLogger that wraps the provided slog.Logger.
func NewSlogLogger(logger *slog.Logger, opts ...Option) (*SlogLogger, error) {
	if logger == nil {
		return nil, ErrNilSlogLogger
	}

	options := DefaultOptions()

	for _, opt := range opts {
		if opt == nil {
			return nil, ErrNilSlogOption
		}

		if err := opt(options); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	return &SlogLogger{
		log: logger,
		opt: options,
	}, nil
}

// Debug sends a debug message to the underlying slog.Logger with the provided context and arguments.
func (sl *SlogLogger) Debug(ctx context.Context, msg string, args ...any) {
	sl.log.DebugContext(ctx, msg, slog.Any(sl.opt.argsKey, args))
}

// Info sends an info message to the underlying slog.Logger with the provided context and arguments.
func (sl *SlogLogger) Info(ctx context.Context, msg string, args ...any) {
	sl.log.InfoContext(ctx, msg, slog.Any(sl.opt.argsKey, args))
}

// Warn sends a warning message to the underlying slog.Logger with the provided context and arguments.
func (sl *SlogLogger) Warn(ctx context.Context, msg string, args ...any) {
	sl.log.WarnContext(ctx, msg, slog.Any(sl.opt.argsKey, args))
}

// Error sends an error message to the underlying slog.Logger with the provided context, error, and arguments.
func (sl *SlogLogger) Error(ctx context.Context, err error, args ...any) {
	sl.log.ErrorContext(ctx, err.Error(), slog.Any(sl.opt.argsKey, args))
}

// Ensures that SlogLogger implements the cakelog.Logger interface.
var _ cakelog.Logger = (*SlogLogger)(nil)
