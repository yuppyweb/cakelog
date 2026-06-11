package decorator

import (
	"context"
	"sync"

	"github.com/yuppyweb/cakelog"
)

// contextLoggerKey is an unexported type used as a unique key for storing
// context values. Using an unexported type prevents collisions with other packages.
type contextLoggerKey struct{}

// WithContextValue returns a new context with the given key-value pair added to its context values.
// If the key already exists, its value is updated.
//
// Note: We copy the sync.Map to ensure isolation between parent and child contexts.
// Without copying, mutations to the child context would affect the parent context as well,
// since context.WithValue only wraps the value without deep copying it.
func WithContextValue(ctx context.Context, key string, value any) context.Context {
	newSm := &sync.Map{}

	if sm, ok := ctx.Value(contextLoggerKey{}).(*sync.Map); ok {
		sm.Range(func(k, v any) bool {
			newSm.Store(k, v)

			return true
		})
	}

	newSm.Store(key, value)

	return context.WithValue(ctx, contextLoggerKey{}, newSm)
}

// ContextValue returns the value associated with the given key in the context,
// or nil if the key is not found.
func ContextValue(ctx context.Context, key string) any {
	sm, ok := ctx.Value(contextLoggerKey{}).(*sync.Map)
	if !ok {
		return nil
	}

	val, ok := sm.Load(key)
	if !ok {
		return nil
	}

	return val
}

// ContextValues returns all key-value pairs stored in the context as a map.
// If no values are stored, it returns an empty map.
func ContextValues(ctx context.Context) map[string]any {
	ctxArgs := make(map[string]any)

	sm, ok := ctx.Value(contextLoggerKey{}).(*sync.Map)
	if !ok {
		return ctxArgs
	}

	sm.Range(func(k, v any) bool {
		if key, ok := k.(string); ok {
			ctxArgs[key] = v
		}

		return true
	})

	return ctxArgs
}

// ContextLogger is a decorator that enriches log messages with context values.
type ContextLogger struct {
	// log is the underlying logger where messages are forwarded.
	log cakelog.Logger
}

// NewContextLogger creates a new ContextLogger that wraps the provided logger.
func NewContextLogger(log cakelog.Logger) (*ContextLogger, error) {
	if log == nil {
		return nil, ErrNilLogger
	}

	return &ContextLogger{log: log}, nil
}

// Debug logs a debug message with context values.
func (cl *ContextLogger) Debug(ctx context.Context, msg string, args ...any) {
	cl.log.Debug(ctx, msg, cl.appendContextArgs(ctx, args)...)
}

// Info logs an info message with context values.
func (cl *ContextLogger) Info(ctx context.Context, msg string, args ...any) {
	cl.log.Info(ctx, msg, cl.appendContextArgs(ctx, args)...)
}

// Warn logs a warning message with context values.
func (cl *ContextLogger) Warn(ctx context.Context, msg string, args ...any) {
	cl.log.Warn(ctx, msg, cl.appendContextArgs(ctx, args)...)
}

// Error logs an error message with context values.
func (cl *ContextLogger) Error(ctx context.Context, err error, args ...any) {
	cl.log.Error(ctx, err, cl.appendContextArgs(ctx, args)...)
}

// appendContextArgs appends context values to the log arguments.
func (cl *ContextLogger) appendContextArgs(ctx context.Context, args []any) []any {
	ctxArgs := ContextValues(ctx)

	if len(ctxArgs) > 0 {
		args = append(args, ctxArgs)
	}

	return args
}

// Ensure that ContextLogger implements the cakelog.Logger interface.
var _ cakelog.Logger = (*ContextLogger)(nil)
