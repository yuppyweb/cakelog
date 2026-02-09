package adapter

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/yuppyweb/cakelog"
)

const DefaultLogrusArgsKey = "context"

type LogrusLogger struct {
	logger  *logrus.Logger
	ArgsKey string
}

func NewLogrusLogger(logger *logrus.Logger) *LogrusLogger {
	return &LogrusLogger{
		logger:  logger,
		ArgsKey: DefaultLogrusArgsKey,
	}
}

func (ll *LogrusLogger) Debug(ctx context.Context, msg string, args ...any) {
	ll.logger.WithField(ll.ArgsKey, args).Debug(msg)
}

func (ll *LogrusLogger) Info(ctx context.Context, msg string, args ...any) {
	ll.logger.WithField(ll.ArgsKey, args).Info(msg)
}

func (ll *LogrusLogger) Warn(ctx context.Context, msg string, args ...any) {
	ll.logger.WithField(ll.ArgsKey, args).Warn(msg)
}

func (ll *LogrusLogger) Error(ctx context.Context, err error, args ...any) {
	ll.logger.WithField(ll.ArgsKey, args).Error(err)
}

var _ cakelog.Logger = (*LogrusLogger)(nil)
