package decorator

import (
	"context"
	"errors"
	"fmt"

	"github.com/yuppyweb/cakelog"
)

// logLevel is the inclusive minimum severity that a levelLogger forwards.
type logLevel int

const (
	// LevelDebug forwards Debug, Info, Warn, and Error.
	LevelDebug logLevel = iota
	// LevelInfo forwards Info, Warn, and Error. Debug is dropped.
	LevelInfo
	// LevelWarn forwards Warn and Error. Debug and Info are dropped.
	LevelWarn
	// LevelError forwards only Error. Debug, Info, and Warn are dropped.
	LevelError
)

// ErrInvalidLogLevel is returned by NewLevel when minLevel is not one of
// LevelDebug, LevelInfo, LevelWarn, or LevelError.
var ErrInvalidLogLevel = errors.New("invalid log level")

// levelLogger is a decorator that drops log calls below minLevel before
// forwarding to the underlying logger.
//
// minLevel is inclusive: a call is forwarded when its level is greater
// than or equal to minLevel. The underlying logger may still drop the
// entry according to its own configuration. Args are not copied. A panic
// in the underlying logger is not recovered and propagates to the caller.
type levelLogger struct {
	// log is the underlying logger where messages that pass the threshold
	// are forwarded.
	log cakelog.Logger

	// minLevel is the inclusive minimum level that is forwarded.
	minLevel logLevel
}

// NewLevel wraps log so each log call below minLevel is dropped and the
// rest are forwarded unchanged.
//
// minLevel is inclusive. Debug is forwarded only at LevelDebug. Info is
// forwarded at LevelDebug and LevelInfo. Warn is forwarded at LevelDebug,
// LevelInfo, and LevelWarn. Error is forwarded at every valid minLevel.
//
// Args are not copied. A panic in the underlying logger is not recovered
// and propagates to the caller. The underlying logger may still drop the
// entry according to its own configuration.
//
// A nil log, including a typed nil such as a nil pointer stored in
// Logger, returns a wrapped ErrNilLogger. A minLevel other than
// LevelDebug, LevelInfo, LevelWarn, or LevelError returns a wrapped
// ErrInvalidLogLevel that includes the numeric value.
func NewLevel(log cakelog.Logger, minLevel logLevel) (cakelog.Logger, error) {
	if err := requireLogger(log); err != nil {
		return nil, fmt.Errorf("level logger: %w", err)
	}

	if minLevel < LevelDebug || minLevel > LevelError {
		return nil, fmt.Errorf("level logger: %w: %d", ErrInvalidLogLevel, minLevel)
	}

	return &levelLogger{log: log, minLevel: minLevel}, nil
}

// Debug forwards the message when minLevel is LevelDebug.
func (ll *levelLogger) Debug(ctx context.Context, msg string, args ...any) {
	if ll.minLevel <= LevelDebug {
		ll.log.Debug(ctx, msg, args...)
	}
}

// Info forwards the message when minLevel is LevelDebug or LevelInfo.
func (ll *levelLogger) Info(ctx context.Context, msg string, args ...any) {
	if ll.minLevel <= LevelInfo {
		ll.log.Info(ctx, msg, args...)
	}
}

// Warn forwards the message when minLevel is LevelDebug, LevelInfo, or LevelWarn.
func (ll *levelLogger) Warn(ctx context.Context, msg string, args ...any) {
	if ll.minLevel <= LevelWarn {
		ll.log.Warn(ctx, msg, args...)
	}
}

// Error forwards the error at every valid minLevel.
func (ll *levelLogger) Error(ctx context.Context, err error, args ...any) {
	if ll.minLevel <= LevelError {
		ll.log.Error(ctx, err, args...)
	}
}

// Ensure levelLogger implements the cakelog.Logger interface.
var _ cakelog.Logger = (*levelLogger)(nil)
