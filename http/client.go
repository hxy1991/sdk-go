package http

import (
	"bytes"
	"context"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/hxy1991/sdk-go/log"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

func Post(ctx context.Context, url string, requestBody []byte) ([]byte, error) {
	var responseBody []byte

	defer func() {
		log.Context(ctx).With("url", url).With("requestBody", string(requestBody), "responseBody", string(responseBody)).Debug()
	}()

	httpClient := xray.Client(&http.Client{Timeout: time.Second * 5})

	request, err := http.NewRequestWithContext(ctx,
		"POST",
		url,
		bytes.NewBuffer(requestBody),
	)
	request.Header.Set("Content-Type", "application/json")

	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer func(body io.ReadCloser) {
		closeErr := body.Close()
		if closeErr != nil {
			log.Context(ctx).Error(closeErr)
		}
	}(response.Body)

	responseBody, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}
