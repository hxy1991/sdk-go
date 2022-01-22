package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type Logger struct {
	_logger *zap.SugaredLogger
}

func defaultLogger() Logger {
	// First, define our level-handling logic.
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})

	// High-priority output should also go to standard error, and low-priority
	// output should also go to standard out.
	stdoutWriteSyncer := zapcore.Lock(os.Stdout)
	stderrWriteSyncer := zapcore.Lock(os.Stderr)

	productionEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	core := zapcore.NewTee(
		zapcore.NewCore(productionEncoder, stderrWriteSyncer, highPriority),
		zapcore.NewCore(productionEncoder, stdoutWriteSyncer, lowPriority),
	)

	// "zap.AddCallerSkip(1)" can locate the real caller because we wrap the zap logger
	_zapLogger := zap.New(core, zap.WithCaller(true), zap.AddStacktrace(zapcore.ErrorLevel), zap.AddCallerSkip(1))
	defer func(zapLogger *zap.Logger) {
		_ = zapLogger.Sync() // flushes buffer, if any
	}(_zapLogger)

	return Logger{_logger: _zapLogger.Sugar()}
}

func (l *Logger) Debug(args ...interface{}) {
	l._logger.Debug(args...)
}

func (l *Logger) Info(args ...interface{}) {
	l._logger.Info(args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l._logger.Warn(args...)
}

func (l *Logger) Error(args ...interface{}) {
	l._logger.Error(args...)
}

func (l *Logger) Warnf(template string, args ...interface{}) {
	l._logger.Warnf(template, args...)
}
