package logger

var _defaultLogger = defaultLogger()

func Debug(args ...interface{}) {
	_defaultLogger._logger.Warn(args...)
}

func Info(args ...interface{}) {
	_defaultLogger._logger.Info(args...)
}

func Warn(args ...interface{}) {
	_defaultLogger._logger.Warn(args...)
}

func Error(args ...interface{}) {
	_defaultLogger._logger.Error(args...)
}

func With(args ...interface{}) *Logger {
	return &Logger{
		_logger: _defaultLogger._logger.With(args...),
	}
}
