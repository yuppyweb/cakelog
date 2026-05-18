package adapter

import (
	"context"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/yuppyweb/cakelog"
)

var (
	ErrNilLogrusLogger = errors.New("is nil logrus.Logger")
	ErrNilLogrusOption = errors.New("is nil logrus option")
)

// LogrusLogger is an adapter that allows using a logrus.Logger as a cakelog.Logger.
type LogrusLogger struct {
	// The underlying logrus.Logger to which log messages will be forwarded.
	log *logrus.Logger

	// The options for configuring the LogrusLogger behavior.
	opt *Options
}

// NewLogrusLogger creates a new LogrusLogger that wraps the provided logrus.Logger.
func NewLogrusLogger(logger *logrus.Logger, opts ...Option) (*LogrusLogger, error) {
	if logger == nil {
		return nil, ErrNilLogrusLogger
	}

	options := DefaultOptions()

	for _, opt := range opts {
		if opt == nil {
			return nil, ErrNilLogrusOption
		}

		if err := opt(options); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	return &LogrusLogger{
		log: logger,
		opt: options,
	}, nil
}

// Debug sends a debug message to the underlying logrus.Logger with the provided context and arguments.
func (ll *LogrusLogger) Debug(ctx context.Context, msg string, args ...any) {
	ll.log.WithContext(ctx).WithField(ll.opt.argsKey, args).Debug(msg)
}

// Info sends an info message to the underlying logrus.Logger with the provided context and arguments.
func (ll *LogrusLogger) Info(ctx context.Context, msg string, args ...any) {
	ll.log.WithContext(ctx).WithField(ll.opt.argsKey, args).Info(msg)
}

// Warn sends a warning message to the underlying logrus.Logger with the provided context and arguments.
func (ll *LogrusLogger) Warn(ctx context.Context, msg string, args ...any) {
	ll.log.WithContext(ctx).WithField(ll.opt.argsKey, args).Warn(msg)
}

// Error sends an error message to the underlying logrus.Logger with the provided context, error, and arguments.
func (ll *LogrusLogger) Error(ctx context.Context, err error, args ...any) {
	ll.log.WithContext(ctx).WithField(ll.opt.argsKey, args).Error(err)
}

// Ensures that LogrusLogger implements the cakelog.Logger interface.
var _ cakelog.Logger = (*LogrusLogger)(nil)
