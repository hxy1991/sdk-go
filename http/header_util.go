package http

import (
	"context"
	"encoding/json"
	"github.com/hxy1991/sdk-go/log"
	"net/http"
	"os"
)

var serverEnv = os.Getenv("SERVER_ENV")

func BuildHeaderForRoute(ctx context.Context, service string, header http.Header) map[string]string {
	headers := map[string]string{}
	headers["Service"] = service
	headers["Env"] = serverEnv

	clientVersion := header.Get("Client-Version")
	if len(clientVersion) == 0 {
		clientDataStr := header.Get("Clientdata")
		clientDataMap := map[string]interface{}{}
		if clientDataStr != "" {
			err := json.Unmarshal([]byte(clientDataStr), &clientDataMap)
			if err != nil {
				log.Context(ctx).Warnf("clientData json.Unmarshal err %+v", err)
			}

			appVersion, ok := clientDataMap["appVersion"]
			if ok {
				clientVersion = appVersion.(string)
			}
		}
	}
	headers["Client-Version"] = clientVersion

	channelId := header.Get("Channel-Id")
	if len(channelId) == 0 {
		// 1-ios, 2-android, 0-其他
		channelId = "0"
	}
	headers["Channel-Id"] = channelId

	return headers
}
