package adapter

import (
	"context"
	"errors"
	"fmt"

	"github.com/yuppyweb/cakelog"
	"go.uber.org/zap"
)

var (
	ErrNilZapLogger = errors.New("is nil zap.Logger")
	ErrNilZapOption = errors.New("is nil zap option")
)

// ZapLogger is an adapter that allows using a zap.Logger as a cakelog.Logger.
type ZapLogger struct {
	// The underlying zap.Logger to which log messages will be forwarded.
	log *zap.Logger

	// The options for configuring the ZapLogger behavior.
	opt *Options
}

// NewZapLogger creates a new ZapLogger that wraps the provided zap.Logger.
func NewZapLogger(logger *zap.Logger, opts ...Option) (*ZapLogger, error) {
	if logger == nil {
		return nil, ErrNilZapLogger
	}

	options := DefaultOptions()

	for _, opt := range opts {
		if opt == nil {
			return nil, ErrNilZapOption
		}

		if err := opt(options); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	return &ZapLogger{
		log: logger,
		opt: options,
	}, nil
}

// Debug sends a debug message to the underlying zap.Logger with the provided context and arguments.
func (zl *ZapLogger) Debug(_ context.Context, msg string, args ...any) {
	zl.log.Debug(msg, zap.Any(zl.opt.argsKey, args))
}

// Info sends an info message to the underlying zap.Logger with the provided context and arguments.
func (zl *ZapLogger) Info(_ context.Context, msg string, args ...any) {
	zl.log.Info(msg, zap.Any(zl.opt.argsKey, args))
}

// Warn sends a warning message to the underlying zap.Logger with the provided context and arguments.
func (zl *ZapLogger) Warn(_ context.Context, msg string, args ...any) {
	zl.log.Warn(msg, zap.Any(zl.opt.argsKey, args))
}

// Error sends an error message to the underlying zap.Logger with the provided context, error, and arguments.
func (zl *ZapLogger) Error(_ context.Context, err error, args ...any) {
	zl.log.Error(err.Error(), zap.Any(zl.opt.argsKey, args))
}

// Ensures that ZapLogger implements the cakelog.Logger interface.
var _ cakelog.Logger = (*ZapLogger)(nil)
