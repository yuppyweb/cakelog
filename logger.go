package cakelog

import (
	"context"
)

type Logger interface {
	Debug(ctx context.Context, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, err error, args ...any)
}

type NopLogger struct{}

func NewNopLogger() *NopLogger {
	return &NopLogger{}
}

func (*NopLogger) Debug(context.Context, string, ...any) {}

func (*NopLogger) Info(context.Context, string, ...any) {}

func (*NopLogger) Warn(context.Context, string, ...any) {}

func (*NopLogger) Error(context.Context, error, ...any) {}

var _ Logger = (*NopLogger)(nil)
