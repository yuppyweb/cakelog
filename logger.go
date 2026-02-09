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

func (*NopLogger) Debug(ctx context.Context, msg string, args ...any) {}

func (*NopLogger) Info(ctx context.Context, msg string, args ...any) {}

func (*NopLogger) Warn(ctx context.Context, msg string, args ...any) {}

func (*NopLogger) Error(ctx context.Context, err error, args ...any) {}

var _ Logger = (*NopLogger)(nil)
