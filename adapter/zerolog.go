package adapter

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/yuppyweb/cakelog"
)

// Is the default key under which context arguments will be stored in zerolog entries.
const DefaultZerologArgsKey = "context"

// Is an adapter that allows using a zerolog.Logger as a cakelog.Logger.
type ZerologLogger struct {
	// The underlying zerolog.Logger to which log messages will be forwarded.
	logger *zerolog.Logger

	// The key under which the context arguments will be stored in zerolog entries.
	ArgsKey string
}

// Creates a new ZerologLogger that wraps the provided zerolog.Logger.
func NewZerologLogger(logger *zerolog.Logger) *ZerologLogger {
	return &ZerologLogger{
		logger:  logger,
		ArgsKey: DefaultZerologArgsKey,
	}
}

// Sends a debug message to the underlying zerolog.Logger with the provided context and arguments.
func (zl *ZerologLogger) Debug(ctx context.Context, msg string, args ...any) {
	zl.logger.Debug().Ctx(ctx).Fields(map[string]any{zl.ArgsKey: args}).Msg(msg)
}

// Sends an info message to the underlying zerolog.Logger with the provided context and arguments.
func (zl *ZerologLogger) Info(ctx context.Context, msg string, args ...any) {
	zl.logger.Info().Ctx(ctx).Fields(map[string]any{zl.ArgsKey: args}).Msg(msg)
}

// Sends a warning message to the underlying zerolog.Logger with the provided context and arguments.
func (zl *ZerologLogger) Warn(ctx context.Context, msg string, args ...any) {
	zl.logger.Warn().Ctx(ctx).Fields(map[string]any{zl.ArgsKey: args}).Msg(msg)
}

// Sends an error message to the underlying zerolog.Logger with the provided context, error, and arguments.
func (zl *ZerologLogger) Error(ctx context.Context, err error, args ...any) {
	zl.logger.Error().Ctx(ctx).Fields(map[string]any{zl.ArgsKey: args}).Msg(err.Error())
}

// Ensures that ZerologLogger implements the cakelog.Logger interface.
var _ cakelog.Logger = (*ZerologLogger)(nil)
