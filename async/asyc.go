package async

import (
	"context"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/hxy1991/sdk-go/constant"
	"github.com/hxy1991/sdk-go/log"
)

func Go(ctx context.Context, name string, fn func(context.Context)) {
	if xray.SdkDisabled() {
		go func() {
			defer func() {
				if e := recover(); e != nil {
					log.Context(ctx).Errorf("%+v", e)
				}
			}()

			// 为什么这里要用context.Background()? 怕外层传入的 context 超时取消了，导致这里的 fn() 也被取消
			newCtx := context.Background()
			newCtx = WithValues(ctx, newCtx)
			fn(newCtx)
		}()
	} else {
		newCtx := xray.DetachContext(ctx)
		newCtx = WithValues(ctx, newCtx)

		xray.CaptureAsync(newCtx, name, func(ctx context.Context) error {
			fn(ctx)
			return nil
		})
	}
}

func WithValues(ctx context.Context, newCtx context.Context) context.Context {
	clientVersion := ctx.Value(constant.ClientVersion)
	if clientVersion != nil {
		newCtx = context.WithValue(newCtx, constant.ClientVersion, clientVersion)
	}

	clientLoginVersion := ctx.Value(constant.ClientLogicVersion)
	if clientLoginVersion != nil {
		newCtx = context.WithValue(newCtx, constant.ClientLogicVersion, clientLoginVersion)
	}

	channelIdKey := ctx.Value(constant.ChannelIdKey)
	if channelIdKey != nil {
		newCtx = context.WithValue(newCtx, constant.ChannelIdKey, channelIdKey)
	}
	return newCtx
}
