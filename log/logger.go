package log

import (
	"context"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/hxy1991/sdk-go/constant"
	"github.com/hxy1991/sdk-go/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type Logger struct {
	_logger *zap.SugaredLogger
	ctx     context.Context
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

	if utils.IsConsoleLog() {
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
	l.withRequestMsg()
	l._logger.Warn(args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.withRequestMsg()
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
	l.withRequestMsg()
	l._logger.Warnf(template, args...)
}

func (l *Logger) Errorf(template string, args ...interface{}) {
	l.withRequestMsg()
	l._logger.Errorf(template, args...)
}

func (l *Logger) withRequestMsg() {
	if l.ctx != nil {
		requestPath, ok := l.ctx.Value(constant.RequestPathKey).(string)
		if ok {
			l._logger = l._logger.With(zap.String("requestPath", requestPath))
		}

		requestBody, ok := l.ctx.Value(constant.RequestBodyKey).(string)
		if ok {
			l._logger = l._logger.With(zap.String("requestBody", requestBody))
		}
	}
}

func (l *Logger) withRequestPath() {
	if l.ctx != nil {
		requestBody, ok := l.ctx.Value(constant.RequestPathKey).(string)
		if ok {
			l._logger = l._logger.With(zap.String("requestPath", requestBody))
		}
	}
}

func (l *Logger) Context(ctx context.Context) (logger *Logger) {
	logger = &Logger{
		_logger: l._logger,
		ctx:     ctx,
	}

	traceId := xray.TraceID(ctx)
	if traceId != "" {
		logger._logger = logger._logger.With(zap.String("xray-trace-id", traceId))
	}

	segment := xray.GetSegment(ctx)
	if segment != nil && segment.ID != "" {
		logger._logger = logger._logger.With(zap.String("xray-segment-id", segment.ID))
	}

	userIdUint64, ok := ctx.Value(constant.UserIdUint64Key).(uint64)
	if ok {
		logger._logger = logger._logger.With(zap.Uint64("userId", userIdUint64))
	}

	accountIdUint64, ok := ctx.Value(constant.AccountIdUint64Key).(uint64)
	if ok {
		logger._logger = logger._logger.With(zap.Uint64("accountId", accountIdUint64))
	}

	serverIdInt, ok := ctx.Value(constant.ServerIdIntKey).(int)
	if ok {
		logger._logger = logger._logger.With(zap.Int("serverId", serverIdInt))
	}

	runVersion := os.Getenv("RUN_VERSION")
	if runVersion != "" {
		logger._logger = logger._logger.With(zap.String("runVersion", runVersion))
	}

	handlerLabelKey, ok := ctx.Value(constant.HandlerLabelKey).(string)
	if ok {
		logger._logger = logger._logger.With(zap.String("handlerLabelKey", handlerLabelKey))
	}

	requestModule, ok := ctx.Value(constant.RequestModuleKey).(string)
	if ok {
		logger._logger = logger._logger.With(zap.String("requestModule", requestModule))
	}

	requestAction, ok := ctx.Value(constant.RequestActionKey).(string)
	if ok {
		logger._logger = logger._logger.With(zap.String("requestAction", requestAction))
	}

	requestSubActions, ok := ctx.Value(constant.RequestSubActionsArrKey).([]string)
	if ok {
		logger._logger = logger._logger.With(zap.Strings("requestSubActions", requestSubActions))
	}

	requestSubActionsMD5, ok := ctx.Value(constant.RequestSubActionsMD5Key).(string)
	if ok {
		logger._logger = logger._logger.With(zap.String("requestSubActionsMD5", requestSubActionsMD5))
	}

	responseErrorCode, ok := ctx.Value(constant.ResponseErrorCodeIntKey).(int)
	if ok {
		logger._logger = logger._logger.With(zap.Int("responseErrorCode", responseErrorCode))
	}

	return logger
}
