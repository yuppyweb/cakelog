package adapter

import (
	"context"

	"github.com/yuppyweb/cakelog"
	"go.uber.org/zap"
)

// Is the default key under which context arguments will be stored in zap entries.
const DefaultZapArgsKey = "context"

// Is an adapter that allows using a zap.Logger as a cakelog.Logger.
type ZapLogger struct {
	// The underlying zap.Logger to which log messages will be forwarded.
	logger *zap.Logger

	// The key under which the context arguments will be stored in zap entries.
	ArgsKey string
}

// Creates a new ZapLogger that wraps the provided zap.Logger.
func NewZapLogger(logger *zap.Logger) *ZapLogger {
	return &ZapLogger{
		logger:  logger,
		ArgsKey: DefaultZapArgsKey,
	}
}

// Sends a debug message to the underlying zap.Logger with the provided context and arguments.
func (zl *ZapLogger) Debug(_ context.Context, msg string, args ...any) {
	zl.logger.Debug(msg, zap.Any(zl.ArgsKey, args))
}

// Sends an info message to the underlying zap.Logger with the provided context and arguments.
func (zl *ZapLogger) Info(_ context.Context, msg string, args ...any) {
	zl.logger.Info(msg, zap.Any(zl.ArgsKey, args))
}

// Sends a warning message to the underlying zap.Logger with the provided context and arguments.
func (zl *ZapLogger) Warn(_ context.Context, msg string, args ...any) {
	zl.logger.Warn(msg, zap.Any(zl.ArgsKey, args))
}

// Sends an error message to the underlying zap.Logger with the provided context, error, and arguments.
func (zl *ZapLogger) Error(_ context.Context, err error, args ...any) {
	zl.logger.Error(err.Error(), zap.Any(zl.ArgsKey, args))
}

// Ensures that ZapLogger implements the cakelog.Logger interface.
var _ cakelog.Logger = (*ZapLogger)(nil)
