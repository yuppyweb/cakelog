package adapter

import (
	"context"

	"github.com/yuppyweb/cakelog"
	"go.uber.org/zap"
)

const DefaultZapArgsKey = "context"

type ZapLogger struct {
	logger  *zap.Logger
	ArgsKey string
}

func NewZapLogger(logger *zap.Logger) *ZapLogger {
	return &ZapLogger{
		logger:  logger,
		ArgsKey: DefaultZapArgsKey,
	}
}

func (zl *ZapLogger) Debug(_ context.Context, msg string, args ...any) {
	zl.logger.Debug(msg, zap.Any(zl.ArgsKey, args))
}

func (zl *ZapLogger) Info(_ context.Context, msg string, args ...any) {
	zl.logger.Info(msg, zap.Any(zl.ArgsKey, args))
}

func (zl *ZapLogger) Warn(_ context.Context, msg string, args ...any) {
	zl.logger.Warn(msg, zap.Any(zl.ArgsKey, args))
}

func (zl *ZapLogger) Error(_ context.Context, err error, args ...any) {
	zl.logger.Error(err.Error(), zap.Any(zl.ArgsKey, args))
}

var _ cakelog.Logger = (*ZapLogger)(nil)
