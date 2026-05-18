package adapter

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/yuppyweb/cakelog"
)

var (
	ErrNilZerologLogger = errors.New("is nil zerolog.Logger")
	ErrNilZerologOption = errors.New("is nil zerolog option")
)

// ZerologLogger is an adapter that allows using a zerolog.Logger as a cakelog.Logger.
type ZerologLogger struct {
	// The underlying zerolog.Logger to which log messages will be forwarded.
	log *zerolog.Logger

	// The options for configuring the ZerologLogger behavior.
	opt *Options
}

// NewZerologLogger creates a new ZerologLogger that wraps the provided zerolog.Logger.
func NewZerologLogger(logger *zerolog.Logger, opts ...Option) (*ZerologLogger, error) {
	if logger == nil {
		return nil, ErrNilZerologLogger
	}

	options := DefaultOptions()

	for _, opt := range opts {
		if opt == nil {
			return nil, ErrNilZerologOption
		}

		if err := opt(options); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	return &ZerologLogger{
		log: logger,
		opt: options,
	}, nil
}

// Debug sends a debug message to the underlying zerolog.Logger with the provided context and arguments.
func (zl *ZerologLogger) Debug(ctx context.Context, msg string, args ...any) {
	zl.log.Debug().Ctx(ctx).Any(zl.opt.argsKey, args).Msg(msg)
}

// Info sends an info message to the underlying zerolog.Logger with the provided context and arguments.
func (zl *ZerologLogger) Info(ctx context.Context, msg string, args ...any) {
	zl.log.Info().Ctx(ctx).Any(zl.opt.argsKey, args).Msg(msg)
}

// Warn sends a warning message to the underlying zerolog.Logger with the provided context and arguments.
func (zl *ZerologLogger) Warn(ctx context.Context, msg string, args ...any) {
	zl.log.Warn().Ctx(ctx).Any(zl.opt.argsKey, args).Msg(msg)
}

// Error sends an error message to the underlying zerolog.Logger with the provided context, error, and arguments.
func (zl *ZerologLogger) Error(ctx context.Context, err error, args ...any) {
	zl.log.Error().Ctx(ctx).Any(zl.opt.argsKey, args).Msg(err.Error())
}

// Ensures that ZerologLogger implements the cakelog.Logger interface.
var _ cakelog.Logger = (*ZerologLogger)(nil)
