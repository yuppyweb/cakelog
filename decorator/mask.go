package decorator

import (
	"context"
	"errors"

	"github.com/yuppyweb/cakelog"
)

var (
	ErrEmptyMaskers = errors.New("is empty maskers")
	ErrNilMasker    = errors.New("is nil masker")
)

// Masker defines the interface for masking log data.
// Each rule can independently mask messages, errors, and arguments.
type Masker interface {
	// MaskMessage applies masking to a log message.
	MaskMessage(msg string) string
	// MaskError applies masking to an error.
	MaskError(err error) error
	// MaskArgument applies masking to a log argument.
	MaskArgument(arg any) any
}

// MaskLogger is a decorator that applies masking rules to all log output.
// It implements the cakelog.Logger interface and processes messages, errors,
// and arguments through a chain of Masker instances.
type MaskLogger struct {
	// log is the underlying logger to which masked log entries will be forwarded.
	log cakelog.Logger

	// maskers is a slice of Masker instances that will be applied sequentially to log data.
	maskers []Masker
}

// NewMaskLogger creates a new MaskLogger that applies the given masking rules
// to all log output in order. Rules are applied sequentially to messages, errors,
// and arguments.
func NewMaskLogger(log cakelog.Logger, maskers ...Masker) (*MaskLogger, error) {
	if log == nil {
		return nil, ErrNilLogger
	}

	if len(maskers) == 0 {
		return nil, ErrEmptyMaskers
	}

	for _, masker := range maskers {
		if masker == nil {
			return nil, ErrNilMasker
		}
	}

	return &MaskLogger{log: log, maskers: maskers}, nil
}

// Debug logs a debug message with masked data.
func (ml *MaskLogger) Debug(ctx context.Context, msg string, args ...any) {
	ml.log.Debug(ctx, ml.maskMessage(msg), ml.maskArguments(args)...)
}

// Info logs an info message with masked data.
func (ml *MaskLogger) Info(ctx context.Context, msg string, args ...any) {
	ml.log.Info(ctx, ml.maskMessage(msg), ml.maskArguments(args)...)
}

// Warn logs a warning message with masked data.
func (ml *MaskLogger) Warn(ctx context.Context, msg string, args ...any) {
	ml.log.Warn(ctx, ml.maskMessage(msg), ml.maskArguments(args)...)
}

// Error logs an error with masked error and arguments.
func (ml *MaskLogger) Error(ctx context.Context, err error, args ...any) {
	ml.log.Error(ctx, ml.maskError(err), ml.maskArguments(args)...)
}

// maskMessage applies all masking rules to a message sequentially.
func (ml *MaskLogger) maskMessage(msg string) string {
	for _, masker := range ml.maskers {
		msg = masker.MaskMessage(msg)
	}

	return msg
}

// maskError applies all masking rules to an error sequentially.
func (ml *MaskLogger) maskError(err error) error {
	for _, masker := range ml.maskers {
		err = masker.MaskError(err)
	}

	return err //nolint:wrapcheck
}

// maskArgument applies all masking rules to a single argument sequentially.
func (ml *MaskLogger) maskArgument(arg any) any {
	for _, masker := range ml.maskers {
		arg = masker.MaskArgument(arg)
	}

	return arg
}

// maskArguments applies masking to all arguments by creating a new masked slice.
func (ml *MaskLogger) maskArguments(args []any) []any {
	if len(args) == 0 {
		return args
	}

	masked := make([]any, len(args))

	for i, arg := range args {
		masked[i] = ml.maskArgument(arg)
	}

	return masked
}

// Ensure MaskLogger implements the cakelog.Logger interface.
var _ cakelog.Logger = (*MaskLogger)(nil)
