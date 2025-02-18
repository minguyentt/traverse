package zlogger

import "go.uber.org/zap"

type zlogger struct {
    *zap.SugaredLogger
}

// Development Logger
func NewLogger() *zlogger {
    return &zlogger{zap.Must(zap.NewDevelopment()).Sugar()}
}

func (l *zlogger) WithArea(area string) *zap.SugaredLogger {
    return l.Named(area)
}
