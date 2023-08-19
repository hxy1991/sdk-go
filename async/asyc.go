package async

import (
	"context"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/hxy1991/sdk-go/constant"
	"github.com/hxy1991/sdk-go/log"
)

func Go(c context.Context, name string, fn func(context.Context)) {
	if xray.SdkDisabled() {
		go func() {
			defer func() {
				if e := recover(); e != nil {
					log.Error(e)
				}
			}()

			fn(context.Background())
		}()
	} else {
		newCtx := xray.DetachContext(c)

		clientVersion := c.Value(constant.ClientVersion)
		if clientVersion != nil {
			newCtx = context.WithValue(newCtx, constant.ClientVersion, clientVersion)
		}

		channelIdKey := c.Value(constant.ChannelIdKey)
		if channelIdKey != nil {
			newCtx = context.WithValue(newCtx, constant.ChannelIdKey, channelIdKey)
		}

		xray.CaptureAsync(newCtx, name, func(ctx context.Context) error {
			fn(ctx)
			return nil
		})
	}
}
