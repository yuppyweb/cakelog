package logrusadapter

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/yuppyweb/cakelog"
	"github.com/yuppyweb/cakelog/adapter"
)

var ErrNilLogrusLogger = errors.New("logrus logger is nil")

// logrusAdapter is an adapter that allows using a logrus.Logger as a cakelog.Logger.
type logrusAdapter struct {
	// The underlying logrus.Logger to which log messages will be forwarded.
	log *logrus.Logger
}

// New creates a new logrusAdapter that implements the cakelog.Logger interface.
func New(logger *logrus.Logger) (cakelog.Logger, error) {
	if logger == nil {
		return nil, ErrNilLogrusLogger
	}

	return &logrusAdapter{log: logger}, nil
}

// Debug sends a debug message to the underlying logrus.Logger with the provided context and arguments.
func (ll *logrusAdapter) Debug(ctx context.Context, msg string, args ...any) {
	ll.log.WithContext(ctx).WithFields(logrusFields(args)).Debug(msg)
}

// Info sends an info message to the underlying logrus.Logger with the provided context and arguments.
func (ll *logrusAdapter) Info(ctx context.Context, msg string, args ...any) {
	ll.log.WithContext(ctx).WithFields(logrusFields(args)).Info(msg)
}

// Warn sends a warning message to the underlying logrus.Logger with the provided context and arguments.
func (ll *logrusAdapter) Warn(ctx context.Context, msg string, args ...any) {
	ll.log.WithContext(ctx).WithFields(logrusFields(args)).Warn(msg)
}

// Error sends an error message to the underlying logrus.Logger with the provided context, error, and arguments.
func (ll *logrusAdapter) Error(ctx context.Context, err error, args ...any) {
	ll.log.WithContext(ctx).WithFields(logrusFields(args)).Error(err)
}

func logrusFields(args []any) logrus.Fields {
	fields := adapter.Fields(args)
	out := make(logrus.Fields, len(fields))

	for _, field := range fields {
		out[field.Key] = field.Value
	}

	return out
}

// Ensures that LogrusLogger implements the cakelog.Logger interface.
var _ cakelog.Logger = (*logrusAdapter)(nil)
