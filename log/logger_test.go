package log

import (
	"context"
	"github.com/hxy1991/sdk-go/constant"
	"testing"
)

func TestNewLogger(t *testing.T) {

}

func TestLogger_Errorf(t *testing.T) {
	ctx := context.WithValue(context.Background(), constant.RequestBodyKey, "123")

	_defaultLogger.Context(ctx).Debugf("test")
	_defaultLogger.Context(ctx).Infof("test")
	_defaultLogger.Context(ctx).Warnf("test")
	_defaultLogger.Context(ctx).Errorf("test")

	_defaultLogger.Context(ctx).Debug("test")
	_defaultLogger.Context(ctx).Info("test")
	_defaultLogger.Context(ctx).Warn("test")
	_defaultLogger.Context(ctx).Error("test")
}
