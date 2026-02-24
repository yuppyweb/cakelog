package decorator

import (
	"context"
	"sync"

	"github.com/yuppyweb/cakelog"
)

type ContextEnricher interface {
	PutContext(ctx context.Context, key string, value any) context.Context
}

type contextLoggerKey struct{}

type ContextLogger struct {
	log cakelog.Logger
	key *contextLoggerKey
}

func NewContextLogger(log cakelog.Logger) *ContextLogger {
	return &ContextLogger{
		log: log,
		key: &contextLoggerKey{},
	}
}

func (cl *ContextLogger) Debug(ctx context.Context, msg string, args ...any) {
	cl.log.Debug(ctx, msg, cl.enrichArgsContext(ctx, args)...)
}

func (cl *ContextLogger) Info(ctx context.Context, msg string, args ...any) {
	cl.log.Info(ctx, msg, cl.enrichArgsContext(ctx, args)...)
}

func (cl *ContextLogger) Warn(ctx context.Context, msg string, args ...any) {
	cl.log.Warn(ctx, msg, cl.enrichArgsContext(ctx, args)...)
}

func (cl *ContextLogger) Error(ctx context.Context, err error, args ...any) {
	cl.log.Error(ctx, err, cl.enrichArgsContext(ctx, args)...)
}

func (cl *ContextLogger) PutContext(ctx context.Context, key string, value any) context.Context {
	newSm := &sync.Map{}

	if sm, ok := ctx.Value(cl.key).(*sync.Map); ok {
		sm.Range(func(k, v any) bool {
			newSm.Store(k, v)

			return true
		})
	}

	newSm.Store(key, value)

	return context.WithValue(ctx, cl.key, newSm)
}

func (cl *ContextLogger) enrichArgsContext(ctx context.Context, args []any) []any {
	if sm, ok := ctx.Value(cl.key).(*sync.Map); ok {
		ctxArgs := make(map[any]any)

		sm.Range(func(k, v any) bool {
			ctxArgs[k] = v

			return true
		})

		if len(ctxArgs) > 0 {
			args = append(args, ctxArgs)
		}
	}

	return args
}

var _ cakelog.Logger = (*ContextLogger)(nil)
var _ ContextEnricher = (*ContextLogger)(nil)
