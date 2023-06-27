package log

import (
	"context"
)

var _defaultLogger = defaultLogger()

var _logger = _defaultLogger._logger

func Debug(args ...interface{}) {
	_logger.Debug(args...)
}

func Info(args ...interface{}) {
	_logger.Info(args...)
}

func Warn(args ...interface{}) {
	_logger.Warn(args...)
}

func Error(args ...interface{}) {
	_logger.Error(args...)
}

func Debugf(template string, args ...interface{}) {
	_logger.Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	_logger.Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	_logger.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	_defaultLogger.withRequestBody()

	_logger.Errorf(template, args...)
}

func With(args ...interface{}) *Logger {
	return _defaultLogger.With(args...)
}

func WithMap(content map[string]interface{}) *Logger {
	return _defaultLogger.WithMap(content)
}

func Context(ctx context.Context) (logger *Logger) {
	return _defaultLogger.Context(ctx)
}
