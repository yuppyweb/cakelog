package decorator

import (
	"context"
	"errors"
	"fmt"

	"github.com/yuppyweb/cakelog"
)

var (
	ErrNilCallbackLogger = errors.New("is nil callback logger")
	ErrNilCallbackOpts   = errors.New("is nil callback options")
	ErrNilCallbackOpt    = errors.New("is nil callback option")
	ErrNilDebugFunc      = errors.New("is nil debug callback function")
	ErrNilInfoFunc       = errors.New("is nil info callback function")
	ErrNilWarnFunc       = errors.New("is nil warn callback function")
	ErrNilErrorFunc      = errors.New("is nil error callback function")
)

// CallbackFunc is a function type that is called when a log event occurs.
type CallbackFunc func(ctx context.Context)

// CallbackOption is a function type used to configure CallbackOptions.
type CallbackOption func(*CallbackOptions) error

// CallbackOptions holds callback functions that are executed on each log level.
type CallbackOptions struct {
	debugFunc CallbackFunc
	infoFunc  CallbackFunc
	warnFunc  CallbackFunc
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
			return ErrNilCallbackOpts
		}

		if fn == nil {
			return ErrNilDebugFunc
		}

		opts.debugFunc = fn

		return nil
	}
}

// WithInfoCallback returns a CallbackOption that sets the info callback function.
func WithInfoCallback(fn CallbackFunc) CallbackOption {
	return func(opts *CallbackOptions) error {
		if opts == nil {
			return ErrNilCallbackOpts
		}

		if fn == nil {
			return ErrNilInfoFunc
		}

		opts.infoFunc = fn

		return nil
	}
}

// WithWarnCallback returns a CallbackOption that sets the warn callback function.
func WithWarnCallback(fn CallbackFunc) CallbackOption {
	return func(opts *CallbackOptions) error {
		if opts == nil {
			return ErrNilCallbackOpts
		}

		if fn == nil {
			return ErrNilWarnFunc
		}

		opts.warnFunc = fn

		return nil
	}
}

// WithErrorCallback returns a CallbackOption that sets the error callback function.
func WithErrorCallback(fn CallbackFunc) CallbackOption {
	return func(opts *CallbackOptions) error {
		if opts == nil {
			return ErrNilCallbackOpts
		}

		if fn == nil {
			return ErrNilErrorFunc
		}

		opts.errorFunc = fn

		return nil
	}
}

// CallbackLogger wraps a logger and executes callbacks on each log level.
type CallbackLogger struct {
	log cakelog.Logger
	opt *CallbackOptions
}

// NewCallbackLogger creates a new CallbackLogger with the given logger and options.
func NewCallbackLogger(log cakelog.Logger, opts ...CallbackOption) (*CallbackLogger, error) {
	if log == nil {
		return nil, ErrNilCallbackLogger
	}

	options := DefaultCallbackOptions()

	for _, opt := range opts {
		if opt == nil {
			return nil, ErrNilCallbackOpt
		}

		if err := opt(options); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	return &CallbackLogger{
		log: log,
		opt: options,
	}, nil
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
