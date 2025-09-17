package logger

import (
	"context"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	logger *logrus.Logger
}

// ContextKey is used for context-based logging
type ContextKey string

const (
	RequestIDKey ContextKey = "request_id"
	UserIDKey    ContextKey = "user_id"
	TraceIDKey   ContextKey = "trace_id"
)

func NewLogger(logLevel string) *Logger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
		},
	})
	level, err := logrus.ParseLevel(strings.ToLower(logLevel))
	if err != nil {
		level = logrus.InfoLevel
		logger.WithError(err).Warn("Invalid log level provided, defaulting to info")
	}
	logger.SetLevel(level)

	return &Logger{
		logger: logger,
	}
}

func (l *Logger) GetLevel() logrus.Level {
	return l.logger.Level
}

func (l *Logger) GetLogger() *logrus.Logger {
	return l.logger
}

func (l *Logger) Trace(args ...interface{}) {
	l.logger.Trace(args...)
}

func (l *Logger) Tracef(format string, args ...interface{}) {
	l.logger.Tracef(format, args...)
}

func (l *Logger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

func (l *Logger) Panic(args ...interface{}) {
	l.logger.Panic(args...)
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	l.logger.Panicf(format, args...)
}

func (l *Logger) WithField(key string, value interface{}) *logrus.Entry {
	return l.logger.WithField(key, value)
}

func (l *Logger) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.logger.WithFields(fields)
}

func (l *Logger) WithError(err error) *logrus.Entry {
	return l.logger.WithError(err)
}

// WithContext creates a logger entry with context values
func (l *Logger) WithContext(ctx context.Context) *logrus.Entry {
	entry := l.logger.WithFields(logrus.Fields{})

	if requestID := ctx.Value(RequestIDKey); requestID != nil {
		entry = entry.WithField("request_id", requestID)
	}

	if userID := ctx.Value(UserIDKey); userID != nil {
		entry = entry.WithField("user_id", userID)
	}

	if traceID := ctx.Value(TraceIDKey); traceID != nil {
		entry = entry.WithField("trace_id", traceID)
	}

	return entry
}

// Context-aware logging methods
func (l *Logger) InfoWithContext(ctx context.Context, args ...interface{}) {
	l.WithContext(ctx).Info(args...)
}

func (l *Logger) ErrorWithContext(ctx context.Context, args ...interface{}) {
	l.WithContext(ctx).Error(args...)
}

func (l *Logger) WarnWithContext(ctx context.Context, args ...interface{}) {
	l.WithContext(ctx).Warn(args...)
}

func (l *Logger) DebugWithContext(ctx context.Context, args ...interface{}) {
	l.WithContext(ctx).Debug(args...)
}
