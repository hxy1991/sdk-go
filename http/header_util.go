package http

import (
	"context"
	"encoding/json"
	"github.com/hxy1991/sdk-go/log"
	"os"
)

var serverEnv = os.Getenv("SERVER_ENV")

func BuildHeaderForRoute(ctx context.Context, service, clientDataStr, channelId string) map[string]string {
	headers := map[string]string{}
	headers["Service"] = service
	headers["Env"] = serverEnv

	clientDataMap := map[string]interface{}{}
	if clientDataStr != "" {
		err := json.Unmarshal([]byte(clientDataStr), &clientDataMap)
		if err != nil {
			log.Context(ctx).Warnf("clientData json.Unmarshal err %+v", err)
		}
	}

	appVersion, ok := clientDataMap["appVersion"]
	if ok {
		headers["Client-Version"] = appVersion.(string)
	} else {
		log.Context(ctx).Warn("appVersion is not exist")
	}

	headers["Channel-Id"] = channelId

	return headers
}
