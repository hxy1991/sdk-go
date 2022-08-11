package log

import (
	"context"
	awsxray "github.com/hxy1991/sdk-go/aws/xray"
	"github.com/hxy1991/sdk-go/utils"
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

	var core zapcore.Core

	if utils.IsDevelopment() {
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		core = zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, stderrWriteSyncer, highPriority),
			zapcore.NewCore(consoleEncoder, stdoutWriteSyncer, lowPriority),
		)
	} else {
		productionEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
		core = zapcore.NewTee(
			zapcore.NewCore(productionEncoder, stderrWriteSyncer, highPriority),
			zapcore.NewCore(productionEncoder, stdoutWriteSyncer, lowPriority),
		)
	}

	hostname, _ := os.Hostname()
	option := zap.Fields(
		zap.String("ip", utils.GetLocalIP()),
		zap.String("instanceId", hostname),
	)

	// "zap.AddCallerSkip(1)" can locate the real caller because we wrap the zap logger
	_zapLogger := zap.New(core, zap.WithCaller(true), zap.AddStacktrace(zapcore.ErrorLevel), zap.AddCallerSkip(1), option)
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

func (l *Logger) With(args ...interface{}) *Logger {
	return &Logger{
		_logger: l._logger.With(args...),
	}
}

func (l *Logger) WithMap(content map[string]interface{}) *Logger {
	list := utils.Map2Slice(content)
	return &Logger{
		_logger: l._logger.With(list...),
	}
}

func (l *Logger) Debugf(template string, args ...interface{}) {
	l._logger.Debugf(template, args...)
}

func (l *Logger) Infof(template string, args ...interface{}) {
	l._logger.Infof(template, args...)
}

func (l *Logger) Warnf(template string, args ...interface{}) {
	l._logger.Warnf(template, args...)
}

func (l *Logger) Errorf(template string, args ...interface{}) {
	l._logger.Errorf(template, args...)
}

func (l *Logger) Context(ctx context.Context) (logger *Logger) {
	logger = &Logger{
		_logger: l._logger,
	}

	segmentId, traceId := awsxray.GetSegmentAndTraceId(ctx)

	if traceId != "" {
		logger._logger = logger._logger.With(zap.String("xray-trace-id", traceId))
	}

	if segmentId != "" {
		logger._logger = logger._logger.With(zap.String("xray-segment-id", segmentId))
	}

	return logger
}
