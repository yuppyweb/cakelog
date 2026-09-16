package decorator

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/yuppyweb/cakelog"
)

var (
	// ErrEmptyMaskers is returned by NewMask when maskers is empty.
	ErrEmptyMaskers = errors.New("maskers are empty")
	// ErrNilMasker is returned by NewMask when a masker is nil.
	ErrNilMasker = errors.New("masker is nil")
)

// Masker masks a log message, error, and arguments.
// Maskers are applied sequentially: the output of one becomes the input of
// the next.
//
// A panic in a Masker method is not recovered and propagates to the caller.
// Implementations must not mutate values they receive, including maps and
// slices inside args; return new values instead. NewMask copies the
// argument slice, not the values it contains.
type Masker interface {
	// MaskMessage applies masking to a log message.
	MaskMessage(msg string) string
	// MaskError applies masking to an error.
	// err may be nil, and the implementation may also return nil.
	MaskError(err error) error
	// MaskArguments applies masking to log arguments.
	// The returned slice may have any length, including zero, and may be nil.
	// A nil result is forwarded to the next masker and the underlying logger
	// as zero arguments.
	//
	// Do not mutate args or nested maps and slices. NewMask copies the
	// argument slice, not the values it contains.
	MaskArguments(args ...any) []any
}

// maskLogger is a decorator that applies maskers to all log output
// before forwarding to the underlying logger.
//
// Maskers are applied sequentially to messages, errors, and arguments.
// The args slice is copied, not the values it contains. A panic in a
// masker is not recovered and propagates to the caller.
type maskLogger struct {
	// log is the underlying logger to which masked log entries will be forwarded.
	log cakelog.Logger

	// maskers is the ordered list of Masker instances applied to each call.
	maskers []Masker
}

// NewMask wraps log so each log call is processed by maskers in order
// before it is forwarded. Maskers are applied sequentially to messages,
// errors, and arguments. The maskers slice is copied; later mutations of
// the caller's slice do not affect the logger.
//
// The args slice is copied for the maskers, not the values it contains.
// A nil call-site args slice is passed to the first masker as nil.
// Do not mutate nested maps and slices; return new values instead.
// A panic in a masker is not recovered and propagates to the caller.
//
// A nil log, including a typed nil such as a nil pointer stored in
// Logger, returns a wrapped ErrNilLogger. An empty maskers list returns
// a wrapped ErrEmptyMaskers. A nil Masker, including a typed nil such as
// a nil pointer stored in Masker, returns a wrapped ErrNilMasker for
// that 1-based position.
func NewMask(log cakelog.Logger, maskers ...Masker) (cakelog.Logger, error) {
	if err := requireLogger(log); err != nil {
		return nil, fmt.Errorf("mask logger: %w", err)
	}

	if len(maskers) == 0 {
		return nil, fmt.Errorf("mask logger: %w", ErrEmptyMaskers)
	}

	maskers = append([]Masker{}, maskers...)

	for i, masker := range maskers {
		if err := requireMasker(masker); err != nil {
			return nil, fmt.Errorf("mask logger: %d: %w", i+1, err)
		}
	}

	return &maskLogger{log: log, maskers: maskers}, nil
}

// Debug masks the message and arguments, then forwards to the underlying logger.
func (ml *maskLogger) Debug(ctx context.Context, msg string, args ...any) {
	ml.log.Debug(ctx, ml.maskMessage(msg), ml.maskArguments(args)...)
}

// Info masks the message and arguments, then forwards to the underlying logger.
func (ml *maskLogger) Info(ctx context.Context, msg string, args ...any) {
	ml.log.Info(ctx, ml.maskMessage(msg), ml.maskArguments(args)...)
}

// Warn masks the message and arguments, then forwards to the underlying logger.
func (ml *maskLogger) Warn(ctx context.Context, msg string, args ...any) {
	ml.log.Warn(ctx, ml.maskMessage(msg), ml.maskArguments(args)...)
}

// Error masks err and the arguments, then forwards to the underlying logger.
// err may be nil; a masker may also return nil.
func (ml *maskLogger) Error(ctx context.Context, err error, args ...any) {
	ml.log.Error(ctx, ml.maskError(err), ml.maskArguments(args)...)
}

// maskMessage applies all maskers to a message sequentially.
func (ml *maskLogger) maskMessage(msg string) string {
	for _, masker := range ml.maskers {
		msg = masker.MaskMessage(msg)
	}

	return msg
}

// maskError applies all maskers to an error sequentially.
// A nil err is passed through; a masker may return nil.
func (ml *maskLogger) maskError(err error) error {
	for _, masker := range ml.maskers {
		err = masker.MaskError(err)
	}

	return err //nolint:wrapcheck
}

// maskArguments copies the args slice (not nested maps or slices) and applies
// all maskers sequentially. A nil args slice is preserved for the first
// masker. A masker may return a slice of any length, including nil.
func (ml *maskLogger) maskArguments(args []any) []any {
	masked := slices.Clone(args)

	for _, masker := range ml.maskers {
		masked = masker.MaskArguments(masked...)
	}

	return masked
}

func requireMasker(masker Masker) error {
	if isNil(masker) {
		return ErrNilMasker
	}

	return nil
}

// Ensure maskLogger implements the cakelog.Logger interface.
var _ cakelog.Logger = (*maskLogger)(nil)
