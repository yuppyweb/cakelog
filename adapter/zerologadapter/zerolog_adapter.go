package zerologadapter

import (
	"context"
	"errors"

	"github.com/rs/zerolog"
	"github.com/yuppyweb/cakelog"
	"github.com/yuppyweb/cakelog/adapter"
)

var ErrNilZerologLogger = errors.New("zerolog logger is nil")

// zerologAdapter is an adapter that allows using a zerolog.Logger as a cakelog.Logger.
type zerologAdapter struct {
	// The underlying zerolog.Logger to which log messages will be forwarded.
	log *zerolog.Logger
}

// New creates a new zerologAdapter that implements the cakelog.Logger interface.
func New(logger *zerolog.Logger) (cakelog.Logger, error) {
	if logger == nil {
		return nil, ErrNilZerologLogger
	}

	return &zerologAdapter{log: logger}, nil
}

// Debug sends a debug message to the underlying zerolog.Logger with the provided context and arguments.
func (zl *zerologAdapter) Debug(ctx context.Context, msg string, args ...any) {
	zerologEvent(zl.log.Debug().Ctx(ctx), args).Msg(msg)
}

// Info sends an info message to the underlying zerolog.Logger with the provided context and arguments.
func (zl *zerologAdapter) Info(ctx context.Context, msg string, args ...any) {
	zerologEvent(zl.log.Info().Ctx(ctx), args).Msg(msg)
}

// Warn sends a warning message to the underlying zerolog.Logger with the provided context and arguments.
func (zl *zerologAdapter) Warn(ctx context.Context, msg string, args ...any) {
	zerologEvent(zl.log.Warn().Ctx(ctx), args).Msg(msg)
}

// Error sends an error message to the underlying zerolog.Logger with the provided context, error, and arguments.
func (zl *zerologAdapter) Error(ctx context.Context, err error, args ...any) {
	msg := ""
	if err != nil {
		msg = err.Error()
	}

	zerologEvent(zl.log.Error().Ctx(ctx), args).Msg(msg)
}

func zerologEvent(event *zerolog.Event, args []any) *zerolog.Event {
	for _, field := range adapter.Fields(args) {
		event = event.Any(field.Key, field.Value)
	}

	return event
}

// Ensures that ZerologLogger implements the cakelog.Logger interface.
var _ cakelog.Logger = (*zerologAdapter)(nil)
