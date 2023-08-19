package http

import (
	"context"
	"github.com/hxy1991/sdk-go/log"
	"os"
)

var serverEnv = os.Getenv("SERVER_ENV")

func BuildHeaderForRoute(ctx context.Context, clientData map[string]interface{}, service string) map[string]string {
	headers := map[string]string{}
	headers["Service"] = service
	headers["Env"] = serverEnv

	appVersion, ok := clientData["appVersion"]
	if ok {
		headers["Client-Version"] = appVersion.(string)
	} else {
		log.Context(ctx).Warn("appVersion is not exist")
	}

	channelId, ok := clientData["Channel-Id"]
	if ok {
		headers["Channel-Id"] = channelId.(string)
	} else {
		log.Context(ctx).Warn("channelId is not exist")
	}

	return headers
}
