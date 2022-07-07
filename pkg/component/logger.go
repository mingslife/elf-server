package component

import (
	"context"
	"fmt"

	"github.com/mingslife/bone"
	"go.uber.org/zap"

	"elf-server/pkg/conf"
)

const (
	LogLevelKey      = "level"
	LogTimestampKey  = "ts"
	LogCallerKey     = "caller"
	LogStacktraceKey = "stacktrace"
	LogMessageKey    = "msg"
	LogTraceKey      = "trace"
)

const (
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
	LogLevelPanic = "panic"
	LogLevelFatal = "fatal"
)

type Logger struct {
	ctx      context.Context
	skip     int
	l        *zap.Logger
	Project  string
	Logstore string
}

func (*Logger) Name() string {
	return "component.logger"
}

func (*Logger) Init() error {
	return nil
}

func (c *Logger) Register() error {
	c.skip = 3

	cfg := conf.GetConfig()
	option := zap.AddCallerSkip(c.skip)

	var (
		logger *zap.Logger
		err    error
	)
	if cfg.Debug {
		logger, err = zap.NewDevelopment(option)
	} else {
		logger, err = zap.NewProduction(option)
	}

	if err != nil {
		return err
	}

	c.l = logger

	return nil
}

func (c *Logger) Unregister() error {
	c.l.Sync()
	return nil
}

func (c *Logger) log(level string, msg any, extra ...map[string]any) {
	// add trace
	if c.ctx != nil {
		traceID := c.ctx.Value(TraceIDKey)
		if len(extra) > 0 {
			extra[0][LogTraceKey] = traceID
		} else {
			extra = []map[string]any{
				{LogTraceKey: traceID},
			}
		}
	}

	msgStr := c.formatMessage(msg)
	kvList, _ := c.formatExtraData(extra)

	c.printLog(level, msgStr, kvList...)
}

func (c *Logger) WithContext(ctx context.Context) bone.Logger {
	newLogger := *c
	newLogger.ctx = ctx
	return &newLogger
}

func (c *Logger) Debug(msg any, extra ...map[string]any) {
	c.log(LogLevelDebug, msg, extra...)
}

func (c *Logger) Info(msg any, extra ...map[string]any) {
	c.log(LogLevelInfo, msg, extra...)
}

func (c *Logger) Warn(msg any, extra ...map[string]any) {
	c.log(LogLevelWarn, msg, extra...)
}

func (c *Logger) Error(msg any, extra ...map[string]any) {
	c.log(LogLevelError, msg, extra...)
}

func (c *Logger) Panic(msg any, extra ...map[string]any) {
	c.log(LogLevelPanic, msg, extra...)
}

func (c *Logger) Fatal(msg any, extra ...map[string]any) {
	c.log(LogLevelFatal, msg, extra...)
}

func (c *Logger) formatMessage(msg any) string {
	return fmt.Sprint(msg)
}

// formatExtraData converts type []map[string]any to type
// []any and map[string]string, but wouldn't handle duplicate keys
func (c *Logger) formatExtraData(extra []map[string]any) ([]any, map[string]string) {
	l, m := []any{}, map[string]string{}
	for _, kv := range extra {
		for k, v := range kv {
			newV := fmt.Sprint(v)

			l = append(l, k)
			l = append(l, newV)

			m[k] = newV
		}
	}
	return l, m
}

func (c *Logger) printLog(level string, msg string, keyAndValues ...any) {
	switch level {
	case LogLevelDebug:
		c.l.Sugar().Debugw(msg, keyAndValues...)
	case LogLevelInfo:
		c.l.Sugar().Infow(msg, keyAndValues...)
	case LogLevelWarn:
		c.l.Sugar().Warnw(msg, keyAndValues...)
	case LogLevelError:
		c.l.Sugar().Errorw(msg, keyAndValues...)
	case LogLevelPanic:
		c.l.Sugar().Panicw(msg, keyAndValues...)
	case LogLevelFatal:
		c.l.Sugar().Fatalw(msg, keyAndValues...)
	}
}

var _ bone.Logger = (*Logger)(nil)
