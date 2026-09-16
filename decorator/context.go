package decorator

import (
	"context"
	"errors"
	"fmt"

	"github.com/yuppyweb/cakelog"
)

// ErrNilContextFields is returned by NewContext when fields is nil.
var ErrNilContextFields = errors.New("context fields is nil")

// contextFields extracts structured fields from ctx for a log call.
// A nil or empty map leaves the call-site arguments unchanged.
// The returned map is prepended as-is; do not mutate it after returning
// if that map may be reused.
type contextFields func(ctx context.Context) map[string]any

// contextLogger is a decorator that prepends fields extracted from ctx
// to each log call as a map[string]any, then forwards to the underlying
// logger.
//
// Call-site args follow the map, so duplicate keys in args override
// extracted fields when an adapter keeps the last value. Wrapping the
// same logger more than once prepends the map on every layer.
//
// A panic in fields is not recovered and propagates to the caller.
type contextLogger struct {
	// log is the underlying logger where messages are forwarded.
	log cakelog.Logger

	// fields extracts a map of fields from the log call's context.
	fields contextFields
}

// NewContext wraps log so each log call prepends fields(ctx) as a
// map[string]any before the call-site arguments. fields is invoked on
// every log call, including when it ignores ctx and returns a static map.
//
// A nil or empty map from fields leaves args unchanged. Call-site args
// follow the map, so duplicate keys in args override extracted fields
// when an adapter keeps the last value. Wrapping the same logger more
// than once prepends the map on every layer.
//
// The map is prepended as-is; fields must not mutate a returned map
// after the call if that map may be reused. A panic in fields is not
// recovered. A nil log, including a typed nil such as a nil pointer
// stored in Logger, returns a wrapped ErrNilLogger.
// A nil fields function returns a wrapped ErrNilContextFields.
func NewContext(log cakelog.Logger, fields contextFields) (cakelog.Logger, error) {
	if err := requireLogger(log); err != nil {
		return nil, fmt.Errorf("context logger: %w", err)
	}

	if fields == nil {
		return nil, fmt.Errorf("context logger: %w", ErrNilContextFields)
	}

	return &contextLogger{log: log, fields: fields}, nil
}

// Debug prepends fields extracted from ctx, then forwards to the underlying logger.
func (cl *contextLogger) Debug(ctx context.Context, msg string, args ...any) {
	cl.log.Debug(ctx, msg, cl.prependContextArgs(ctx, args)...)
}

// Info prepends fields extracted from ctx, then forwards to the underlying logger.
func (cl *contextLogger) Info(ctx context.Context, msg string, args ...any) {
	cl.log.Info(ctx, msg, cl.prependContextArgs(ctx, args)...)
}

// Warn prepends fields extracted from ctx, then forwards to the underlying logger.
func (cl *contextLogger) Warn(ctx context.Context, msg string, args ...any) {
	cl.log.Warn(ctx, msg, cl.prependContextArgs(ctx, args)...)
}

// Error prepends fields extracted from ctx, then forwards to the underlying logger.
func (cl *contextLogger) Error(ctx context.Context, err error, args ...any) {
	cl.log.Error(ctx, err, cl.prependContextArgs(ctx, args)...)
}

// prependContextArgs returns args with the map from fields prepended.
// A nil or empty map leaves args unchanged. The returned slice does not
// alias the caller's argument storage.
func (cl *contextLogger) prependContextArgs(ctx context.Context, args []any) []any {
	fields := cl.fields(ctx)

	if len(fields) == 0 {
		return args
	}

	return append([]any{fields}, args...)
}

// Ensure contextLogger implements the cakelog.Logger interface.
var _ cakelog.Logger = (*contextLogger)(nil)
