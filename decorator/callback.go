package decorator

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/yuppyweb/cakelog"
)

var (
	// ErrNilCallbackOption is returned by NewCallback when an option is nil.
	ErrNilCallbackOption = errors.New("callback option is nil")
	// ErrNilCallbackDebugFunc is returned by NewCallback when WithDebugCallback
	// is given a nil function.
	ErrNilCallbackDebugFunc = errors.New("debug callback is nil")
	// ErrNilCallbackInfoFunc is returned by NewCallback when WithInfoCallback
	// is given a nil function.
	ErrNilCallbackInfoFunc = errors.New("info callback is nil")
	// ErrNilCallbackWarnFunc is returned by NewCallback when WithWarnCallback
	// is given a nil function.
	ErrNilCallbackWarnFunc = errors.New("warn callback is nil")
	// ErrNilCallbackErrorFunc is returned by NewCallback when WithErrorCallback
	// is given a nil function.
	ErrNilCallbackErrorFunc = errors.New("error callback is nil")
)

// msgCallbackFunc is called after Debug, Info, or Warn forwards the entry
// to the underlying logger. It receives the same context, message, and
// arguments as the log call. The args slice is copied; nested maps and
// slices are not.
//
// The callback runs on the method call, even if the underlying logger later
// discards the entry. If the underlying logger panics, the callback does
// not run. A panic in the callback is not recovered.
type msgCallbackFunc func(ctx context.Context, msg string, args ...any)

// errCallbackFunc is called after Error forwards the entry to the underlying
// logger. It receives the same context, error, and arguments as the log call.
// err may be nil. The args slice is copied; nested maps and slices are not.
//
// The callback runs on the method call, even if the underlying logger later
// discards the entry. If the underlying logger panics, the callback does
// not run. A panic in the callback is not recovered.
type errCallbackFunc func(ctx context.Context, err error, args ...any)

// CallbackOption configures a callback for NewCallback. Pass options to
// NewCallback; do not invoke them directly. If the same level is configured
// more than once, the last value wins.
type CallbackOption func(*callbackOptions) error

// callbackOptions holds the callback for each log level.
// Unset callbacks default to no-op functions.
type callbackOptions struct {
	// debugFunc is called after a Debug log call.
	debugFunc msgCallbackFunc

	// infoFunc is called after an Info log call.
	infoFunc msgCallbackFunc

	// warnFunc is called after a Warn log call.
	warnFunc msgCallbackFunc

	// errorFunc is called after an Error log call.
	errorFunc errCallbackFunc
}

// defaultCallbackOptions returns a no-op callback for every log level.
func defaultCallbackOptions() *callbackOptions {
	noopMsg := func(context.Context, string, ...any) {}
	noopErr := func(context.Context, error, ...any) {}

	return &callbackOptions{
		debugFunc: noopMsg,
		infoFunc:  noopMsg,
		warnFunc:  noopMsg,
		errorFunc: noopErr,
	}
}

// WithDebugCallback sets the Debug callback used by NewCallback.
// A nil function returns ErrNilCallbackDebugFunc when NewCallback applies
// the option.
func WithDebugCallback(fn msgCallbackFunc) CallbackOption {
	return func(opts *callbackOptions) error {
		if fn == nil {
			return ErrNilCallbackDebugFunc
		}

		opts.debugFunc = fn

		return nil
	}
}

// WithInfoCallback sets the Info callback used by NewCallback.
// A nil function returns ErrNilCallbackInfoFunc when NewCallback applies
// the option.
func WithInfoCallback(fn msgCallbackFunc) CallbackOption {
	return func(opts *callbackOptions) error {
		if fn == nil {
			return ErrNilCallbackInfoFunc
		}

		opts.infoFunc = fn

		return nil
	}
}

// WithWarnCallback sets the Warn callback used by NewCallback.
// A nil function returns ErrNilCallbackWarnFunc when NewCallback applies
// the option.
func WithWarnCallback(fn msgCallbackFunc) CallbackOption {
	return func(opts *callbackOptions) error {
		if fn == nil {
			return ErrNilCallbackWarnFunc
		}

		opts.warnFunc = fn

		return nil
	}
}

// WithErrorCallback sets the Error callback used by NewCallback.
// A nil function returns ErrNilCallbackErrorFunc when NewCallback applies
// the option.
func WithErrorCallback(fn errCallbackFunc) CallbackOption {
	return func(opts *callbackOptions) error {
		if fn == nil {
			return ErrNilCallbackErrorFunc
		}

		opts.errorFunc = fn

		return nil
	}
}

// callbackLogger is a decorator that runs a per-level callback after
// forwarding each log call to the underlying logger.
//
// Callbacks fire on the Logger method call, not on actual emission by the
// adapter. The same ctx, message or error, and args are passed to the
// callback. The args slice is copied for the callback, not the values it
// contains; the underlying logger receives the caller's slice. A panic in
// the underlying logger is not recovered and skips the callback. A panic
// in a callback is not recovered and propagates to the caller.
type callbackLogger struct {
	// log is the underlying logger; callbacks run after each call is forwarded.
	log cakelog.Logger

	// opt holds the callback functions for each log level.
	opt *callbackOptions
}

// NewCallback wraps log with per-level callbacks. Omitted levels use no-op
// callbacks. If the same level is configured more than once, the last
// option wins.
//
// Callbacks run after the call is forwarded to the underlying logger, on
// the Logger method call rather than on actual emission by the adapter.
// The underlying logger may still drop the entry. If the underlying logger
// panics, the callback does not run. A panic in a callback is not recovered
// and propagates to the caller.
//
// The same ctx, message or error, and args are passed to the callback.
// The args slice is copied for the callback, not the values it contains;
// the underlying logger receives the caller's slice. Do not mutate nested
// maps and slices.
//
// A nil log, including a typed nil such as a nil pointer stored in Logger,
// returns a wrapped ErrNilLogger. A nil option returns a wrapped
// ErrNilCallbackOption. A nil callback function returns a wrapped
// ErrNilCallbackDebugFunc, ErrNilCallbackInfoFunc, ErrNilCallbackWarnFunc,
// or ErrNilCallbackErrorFunc.
func NewCallback(log cakelog.Logger, opts ...CallbackOption) (cakelog.Logger, error) {
	if err := requireLogger(log); err != nil {
		return nil, fmt.Errorf("callback logger: %w", err)
	}

	options := defaultCallbackOptions()

	for _, opt := range opts {
		if opt == nil {
			return nil, fmt.Errorf("callback logger: %w", ErrNilCallbackOption)
		}

		if err := opt(options); err != nil {
			return nil, fmt.Errorf("callback logger: %w", err)
		}
	}

	return &callbackLogger{log: log, opt: options}, nil
}

// Debug forwards the message to the underlying logger, then calls the debug callback.
func (cl *callbackLogger) Debug(ctx context.Context, msg string, args ...any) {
	cl.log.Debug(ctx, msg, args...)
	cl.opt.debugFunc(ctx, msg, slices.Clone(args)...)
}

// Info forwards the message to the underlying logger, then calls the info callback.
func (cl *callbackLogger) Info(ctx context.Context, msg string, args ...any) {
	cl.log.Info(ctx, msg, args...)
	cl.opt.infoFunc(ctx, msg, slices.Clone(args)...)
}

// Warn forwards the message to the underlying logger, then calls the warn callback.
func (cl *callbackLogger) Warn(ctx context.Context, msg string, args ...any) {
	cl.log.Warn(ctx, msg, args...)
	cl.opt.warnFunc(ctx, msg, slices.Clone(args)...)
}

// Error forwards the error to the underlying logger, then calls the error callback.
func (cl *callbackLogger) Error(ctx context.Context, err error, args ...any) {
	cl.log.Error(ctx, err, args...)
	cl.opt.errorFunc(ctx, err, slices.Clone(args)...)
}

// Ensure callbackLogger implements the cakelog.Logger interface.
var _ cakelog.Logger = (*callbackLogger)(nil)
