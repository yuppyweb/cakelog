package adapter

import (
	"context"
	"log/slog"

	"github.com/yuppyweb/cakelog"
)

const DefaultSlogArgsKey = "context"

type SlogLogger struct {
	Logger  *slog.Logger
	ArgsKey string
}

func NewSlogLogger(logger *slog.Logger) *SlogLogger {
	return &SlogLogger{
		Logger:  logger,
		ArgsKey: DefaultSlogArgsKey,
	}
}

func (sl *SlogLogger) Debug(ctx context.Context, msg string, args ...any) {
	sl.Logger.DebugContext(ctx, msg, slog.Any(sl.ArgsKey, args))
}

func (sl *SlogLogger) Info(ctx context.Context, msg string, args ...any) {
	sl.Logger.InfoContext(ctx, msg, slog.Any(sl.ArgsKey, args))
}

func (sl *SlogLogger) Warn(ctx context.Context, msg string, args ...any) {
	sl.Logger.WarnContext(ctx, msg, slog.Any(sl.ArgsKey, args))
}

func (sl *SlogLogger) Error(ctx context.Context, err error, args ...any) {
	sl.Logger.ErrorContext(ctx, err.Error(), slog.Any(sl.ArgsKey, args))
}

var _ cakelog.Logger = (*SlogLogger)(nil)
