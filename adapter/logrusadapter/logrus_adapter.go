package logrusadapter

import (
	"context"
	"errors"
	"reflect"

	"github.com/sirupsen/logrus"
	"github.com/yuppyweb/cakelog"
	"github.com/yuppyweb/cakelog/adapter"
)

// ErrNilLogrusLogger is returned by New when logger is nil, including a typed
// nil *logrus.Logger or *logrus.Entry stored in FieldLogger.
var ErrNilLogrusLogger = errors.New("logrus logger is nil")

// logrusAdapter wraps a logrus.FieldLogger as a cakelog.Logger.
type logrusAdapter struct {
	// log is the underlying FieldLogger (*logrus.Logger, *logrus.Entry,
	// or another implementation).
	log logrus.FieldLogger
}

// New wraps logger as a cakelog.Logger. logger may be a *logrus.Logger,
// a *logrus.Entry, or any other logrus.FieldLogger.
//
// A nil logger, including a typed nil such as a nil *logrus.Logger or
// *logrus.Entry stored in FieldLogger, returns ErrNilLogrusLogger.
func New(logger logrus.FieldLogger) (cakelog.Logger, error) {
	if isNilLogrusLogger(logger) {
		return nil, ErrNilLogrusLogger
	}

	return &logrusAdapter{log: logger}, nil
}

// Debug sends a debug message to the underlying FieldLogger with the provided context and arguments.
func (ll *logrusAdapter) Debug(ctx context.Context, msg string, args ...any) {
	ll.log.WithFields(logrusFields(args)).WithContext(ctx).Debug(msg)
}

// Info sends an info message to the underlying FieldLogger with the provided context and arguments.
func (ll *logrusAdapter) Info(ctx context.Context, msg string, args ...any) {
	ll.log.WithFields(logrusFields(args)).WithContext(ctx).Info(msg)
}

// Warn sends a warning message to the underlying FieldLogger with the provided context and arguments.
func (ll *logrusAdapter) Warn(ctx context.Context, msg string, args ...any) {
	ll.log.WithFields(logrusFields(args)).WithContext(ctx).Warn(msg)
}

// Error sends an error message to the underlying FieldLogger with the provided context, error, and arguments.
// A non-nil err is also attached with WithError before call-site fields. A later
// field named "error" overwrites it, matching logrus last-wins maps.
func (ll *logrusAdapter) Error(ctx context.Context, err error, args ...any) {
	ll.log.WithError(err).WithFields(logrusFields(args)).WithContext(ctx).Error(err)
}

func logrusFields(args []any) logrus.Fields {
	fields := adapter.Fields(args)
	out := make(logrus.Fields, len(fields))

	for _, field := range fields {
		out[field.Key] = field.Value
	}

	return out
}

func isNilLogrusLogger(logger logrus.FieldLogger) bool {
	if logger == nil {
		return true
	}

	val := reflect.ValueOf(logger)

	//nolint:exhaustive // IsNil is valid only for nillable kinds.
	switch val.Kind() {
	case reflect.Pointer,
		reflect.Interface,
		reflect.Slice,
		reflect.Map,
		reflect.Chan,
		reflect.Func,
		reflect.UnsafePointer:
		if val.IsNil() {
			return true
		}
	}

	return false
}

// Ensures that logrusAdapter implements the cakelog.Logger interface.
var _ cakelog.Logger = (*logrusAdapter)(nil)
