package adapter

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/yuppyweb/cakelog"
)

const DefaultZerologArgsKey = "context"

type ZerologLogger struct {
	logger  *zerolog.Logger
	ArgsKey string
}

func NewZerologLogger(logger *zerolog.Logger) *ZerologLogger {
	return &ZerologLogger{
		logger:  logger,
		ArgsKey: DefaultZerologArgsKey,
	}
}

func (zl *ZerologLogger) Debug(ctx context.Context, msg string, args ...any) {
	zl.logger.Debug().Ctx(ctx).Fields(map[string]any{zl.ArgsKey: args}).Msg(msg)
}

func (zl *ZerologLogger) Info(ctx context.Context, msg string, args ...any) {
	zl.logger.Info().Ctx(ctx).Fields(map[string]any{zl.ArgsKey: args}).Msg(msg)
}

func (zl *ZerologLogger) Warn(ctx context.Context, msg string, args ...any) {
	zl.logger.Warn().Ctx(ctx).Fields(map[string]any{zl.ArgsKey: args}).Msg(msg)
}

func (zl *ZerologLogger) Error(ctx context.Context, err error, args ...any) {
	zl.logger.Error().Ctx(ctx).Fields(map[string]any{zl.ArgsKey: args}).Msg(err.Error())
}

var _ cakelog.Logger = (*ZerologLogger)(nil)
