package async

import (
	"context"
	"github.com/aws/aws-xray-sdk-go/xray"
	"sdk-go/log"
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
		xray.CaptureAsync(newCtx, name, func(ctx context.Context) error {
			fn(ctx)
			return nil
		})
	}
}
