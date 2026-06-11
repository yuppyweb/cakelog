package decorator

import (
	"context"
	"errors"
	"fmt"

	"github.com/yuppyweb/cakelog"
)

var (
	ErrNilCallbackOptions   = errors.New("is nil callback options")
	ErrNilCallbackOption    = errors.New("is nil callback option")
	ErrNilCallbackDebugFunc = errors.New("is nil debug callback function")
	ErrNilCallbackInfoFunc  = errors.New("is nil info callback function")
	ErrNilCallbackWarnFunc  = errors.New("is nil warn callback function")
	ErrNilCallbackErrorFunc = errors.New("is nil error callback function")
)

// CallbackFunc is a function type that is called when a log event occurs.
type CallbackFunc func(ctx context.Context)

// CallbackOption is a function type used to configure CallbackOptions.
type CallbackOption func(*CallbackOptions) error

// CallbackOptions holds callback functions that are executed on each log level.
type CallbackOptions struct {
	// debugFunc is called when a debug log event occurs.
	debugFunc CallbackFunc

	// infoFunc is called when an info log event occurs.
	infoFunc CallbackFunc

	// warnFunc is called when a warning log event occurs.
	warnFunc CallbackFunc

	// errorFunc is called when an error log event occurs.
	errorFunc CallbackFunc
}

// DefaultCallbackOptions returns a new CallbackOptions with default callback functions.
func DefaultCallbackOptions() *CallbackOptions {
	defaultCallbackFunc := func(context.Context) {}

	return &CallbackOptions{
		debugFunc: defaultCallbackFunc,
		infoFunc:  defaultCallbackFunc,
		warnFunc:  defaultCallbackFunc,
		errorFunc: defaultCallbackFunc,
	}
}

// WithDebugCallback returns a CallbackOption that sets the debug callback function.
func WithDebugCallback(fn CallbackFunc) CallbackOption {
	return func(opts *CallbackOptions) error {
		if opts == nil {
			return ErrNilCallbackOptions
		}

		if fn == nil {
			return ErrNilCallbackDebugFunc
		}

		opts.debugFunc = fn

		return nil
	}
}

// WithInfoCallback returns a CallbackOption that sets the info callback function.
func WithInfoCallback(fn CallbackFunc) CallbackOption {
	return func(opts *CallbackOptions) error {
		if opts == nil {
			return ErrNilCallbackOptions
		}

		if fn == nil {
			return ErrNilCallbackInfoFunc
		}

		opts.infoFunc = fn

		return nil
	}
}

// WithWarnCallback returns a CallbackOption that sets the warn callback function.
func WithWarnCallback(fn CallbackFunc) CallbackOption {
	return func(opts *CallbackOptions) error {
		if opts == nil {
			return ErrNilCallbackOptions
		}

		if fn == nil {
			return ErrNilCallbackWarnFunc
		}

		opts.warnFunc = fn

		return nil
	}
}

// WithErrorCallback returns a CallbackOption that sets the error callback function.
func WithErrorCallback(fn CallbackFunc) CallbackOption {
	return func(opts *CallbackOptions) error {
		if opts == nil {
			return ErrNilCallbackOptions
		}

		if fn == nil {
			return ErrNilCallbackErrorFunc
		}

		opts.errorFunc = fn

		return nil
	}
}

// CallbackLogger is a logger decorator that executes callback functions on each log event.
type CallbackLogger struct {
	// log is the underlying logger to which log messages will be forwarded after executing callbacks.
	log cakelog.Logger

	// opt holds the configuration options for the CallbackLogger, such as the callback functions for each log level.
	opt *CallbackOptions
}

// NewCallbackLogger creates a new CallbackLogger with the given logger and options.
func NewCallbackLogger(log cakelog.Logger, opts ...CallbackOption) (*CallbackLogger, error) {
	if log == nil {
		return nil, ErrNilLogger
	}

	options := DefaultCallbackOptions()

	for _, opt := range opts {
		if opt == nil {
			return nil, ErrNilCallbackOption
		}

		if err := opt(options); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	return &CallbackLogger{log: log, opt: options}, nil
}

// Debug logs a debug message and calls the debug callback function.
func (cl *CallbackLogger) Debug(ctx context.Context, msg string, args ...any) {
	cl.log.Debug(ctx, msg, args...)
	cl.opt.debugFunc(ctx)
}

// Info logs an info message and calls the info callback function.
func (cl *CallbackLogger) Info(ctx context.Context, msg string, args ...any) {
	cl.log.Info(ctx, msg, args...)
	cl.opt.infoFunc(ctx)
}

// Warn logs a warning message and calls the warn callback function.
func (cl *CallbackLogger) Warn(ctx context.Context, msg string, args ...any) {
	cl.log.Warn(ctx, msg, args...)
	cl.opt.warnFunc(ctx)
}

// Error logs an error and calls the error callback function.
func (cl *CallbackLogger) Error(ctx context.Context, err error, args ...any) {
	cl.log.Error(ctx, err, args...)
	cl.opt.errorFunc(ctx)
}

// Ensure CallbackLogger implements the cakelog.Logger interface.
var _ cakelog.Logger = (*CallbackLogger)(nil)
