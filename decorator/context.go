package decorator

import (
	"context"
	"sync"

	"github.com/yuppyweb/cakelog"
)

// Is an interface that defines a method for adding key-value pairs to a context,
// which can then be included in log messages.
type ContextEnricher interface {
	// Adds a key-value pair to the context that will be included in all log messages
	// sent through this ContextLogger.
	PutContext(ctx context.Context, key string, value any) context.Context
}

// Is a context key type used to store context values in the context.
// This is unexported to prevent collisions with other context keys.
type contextLoggerKey struct{}

// Is a decorator that enriches log messages with context values stored using the PutContext method.
type ContextLogger struct {
	// The underlying cakelog.Logger to which log messages will be forwarded.
	log cakelog.Logger

	// The key used to store context values in the context.
	// This should be unique to avoid collisions with other context values.
	key *contextLoggerKey
}

// Creates a new ContextLogger that wraps the provided cakelog.Logger.
func NewContextLogger(log cakelog.Logger) *ContextLogger {
	return &ContextLogger{
		log: log,
		key: &contextLoggerKey{},
	}
}

// Sends a debug message to the underlying cakelog.Logger with the provided context and arguments,
// enriched with context values.
func (cl *ContextLogger) Debug(ctx context.Context, msg string, args ...any) {
	cl.log.Debug(ctx, msg, cl.enrichArgsContext(ctx, args)...)
}

// Sends an info message to the underlying cakelog.Logger with the provided context and arguments,
// enriched with context values.
func (cl *ContextLogger) Info(ctx context.Context, msg string, args ...any) {
	cl.log.Info(ctx, msg, cl.enrichArgsContext(ctx, args)...)
}

// Sends a warning message to the underlying cakelog.Logger with the provided context and arguments,
// enriched with context values.
func (cl *ContextLogger) Warn(ctx context.Context, msg string, args ...any) {
	cl.log.Warn(ctx, msg, cl.enrichArgsContext(ctx, args)...)
}

// Sends an error message to the underlying cakelog.Logger with the provided context and arguments,
// enriched with context values.
func (cl *ContextLogger) Error(ctx context.Context, err error, args ...any) {
	cl.log.Error(ctx, err, cl.enrichArgsContext(ctx, args)...)
}

// Puts a key-value pair into the context that will be included in all log messages sent through this ContextLogger.
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

// Helper method to enrich log arguments with context values stored in the context.
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

var (
	// Ensures that ContextLogger implements both the cakelog.Logger and ContextEnricher interfaces.
	_ cakelog.Logger = (*ContextLogger)(nil)

	// Ensures that ContextLogger implements the ContextEnricher interface.
	_ ContextEnricher = (*ContextLogger)(nil)
)
