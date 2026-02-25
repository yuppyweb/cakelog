package adapter

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/yuppyweb/cakelog"
)

// Is the default key under which context arguments will be stored in logrus entries.
const DefaultLogrusArgsKey = "context"

// Is an adapter that allows using a logrus.Logger as a cakelog.Logger.
type LogrusLogger struct {
	// The underlying logrus.Logger to which log messages will be forwarded.
	logger *logrus.Logger

	// The key under which the context arguments will be stored in logrus entries.
	ArgsKey string
}

// Creates a new LogrusLogger that wraps the provided logrus.Logger.
func NewLogrusLogger(logger *logrus.Logger) *LogrusLogger {
	return &LogrusLogger{
		logger:  logger,
		ArgsKey: DefaultLogrusArgsKey,
	}
}

// Sends a debug message to the underlying logrus.Logger with the provided context and arguments.
func (ll *LogrusLogger) Debug(ctx context.Context, msg string, args ...any) {
	ll.logger.WithContext(ctx).WithField(ll.ArgsKey, args).Debug(msg)
}

// Sends an info message to the underlying logrus.Logger with the provided context and arguments.
func (ll *LogrusLogger) Info(ctx context.Context, msg string, args ...any) {
	ll.logger.WithContext(ctx).WithField(ll.ArgsKey, args).Info(msg)
}

// Sends a warning message to the underlying logrus.Logger with the provided context and arguments.
func (ll *LogrusLogger) Warn(ctx context.Context, msg string, args ...any) {
	ll.logger.WithContext(ctx).WithField(ll.ArgsKey, args).Warn(msg)
}

// Sends an error message to the underlying logrus.Logger with the provided context, error, and arguments.
func (ll *LogrusLogger) Error(ctx context.Context, err error, args ...any) {
	ll.logger.WithContext(ctx).WithField(ll.ArgsKey, args).Error(err)
}

// Ensures that LogrusLogger implements the cakelog.Logger interface.
var _ cakelog.Logger = (*LogrusLogger)(nil)
