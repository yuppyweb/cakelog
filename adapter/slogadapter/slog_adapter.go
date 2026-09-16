package slogadapter

import (
	"context"
	"errors"
	"log/slog"

	"github.com/yuppyweb/cakelog"
	"github.com/yuppyweb/cakelog/adapter"
)

// slogArgsPairSize is the number of slice elements used to store one slog attribute.
const slogArgsPairSize = 2

// ErrNilSlogLogger is returned by New when logger is nil.
var ErrNilSlogLogger = errors.New("slog logger is nil")

// slogAdapter is an adapter that allows using a slog.Logger as a cakelog.Logger.
type slogAdapter struct {
	// The underlying slog.Logger to which log messages will be forwarded.
	log *slog.Logger
}

// New creates a new slogAdapter that implements the cakelog.Logger interface.
func New(logger *slog.Logger) (cakelog.Logger, error) {
	if logger == nil {
		return nil, ErrNilSlogLogger
	}

	return &slogAdapter{log: logger}, nil
}

// Debug sends a debug message to the underlying slog.Logger with the provided context and arguments.
func (sl *slogAdapter) Debug(ctx context.Context, msg string, args ...any) {
	sl.log.DebugContext(ctx, msg, flattenArgs(args)...)
}

// Info sends an info message to the underlying slog.Logger with the provided context and arguments.
func (sl *slogAdapter) Info(ctx context.Context, msg string, args ...any) {
	sl.log.InfoContext(ctx, msg, flattenArgs(args)...)
}

// Warn sends a warning message to the underlying slog.Logger with the provided context and arguments.
func (sl *slogAdapter) Warn(ctx context.Context, msg string, args ...any) {
	sl.log.WarnContext(ctx, msg, flattenArgs(args)...)
}

// Error sends an error message to the underlying slog.Logger with the provided context, error, and arguments.
// A non-nil err is also attached as the slog attribute "error" before call-site args.
func (sl *slogAdapter) Error(ctx context.Context, err error, args ...any) {
	msg := ""
	attrs := flattenArgs(args)

	if err != nil {
		msg = err.Error()
		attrs = append([]any{"error", err}, attrs...)
	}

	sl.log.ErrorContext(ctx, msg, attrs...)
}

func flattenArgs(args []any) []any {
	fields := adapter.Fields(args)

	if len(fields) == 0 {
		return nil
	}

	out := make([]any, len(fields)*slogArgsPairSize)

	for i, field := range fields {
		out[i*slogArgsPairSize] = field.Key
		out[i*slogArgsPairSize+1] = field.Value
	}

	return out
}

// Ensures that slogAdapter implements the cakelog.Logger interface.
var _ cakelog.Logger = (*slogAdapter)(nil)
