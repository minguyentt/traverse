package tracer

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/tracelog"
)

type TraceLogger struct {
	logger   *slog.Logger
	levelKey string
}

func NewLogger(l *slog.Logger, levelKey string) *TraceLogger {
	logger := &TraceLogger{
		logger:   l,
		levelKey: levelKey,
	}

	return logger
}

func (l *TraceLogger) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {
	attrs := make([]slog.Attr, 0, len(data))

	for k, v := range data {
		attrs = append(attrs, slog.Any(k, v))
	}

	var lv slog.Level

	switch level {
	case tracelog.LogLevelTrace:
		lv = slog.LevelDebug - 1
		attrs = append(attrs, slog.Any("PGX_LOG_LEVEL", level))
	case tracelog.LogLevelDebug:
		lv = slog.LevelDebug
	case tracelog.LogLevelInfo:
		lv = slog.LevelInfo
	case tracelog.LogLevelWarn:
		lv = slog.LevelWarn
	case tracelog.LogLevelError:
		lv = slog.LevelError
	default:
		lv = slog.LevelError
		attrs = append(
			attrs,
			slog.Any(l.levelKey, fmt.Errorf("error pgx log level: %v", level)),
		)
	}
	l.logger.LogAttrs(ctx, lv, msg, attrs...)
}
