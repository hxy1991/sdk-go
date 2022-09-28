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

func Send(ctx context.Context, url, method string, requestBody []byte, headers map[string]string) (int, []byte, error) {
	var responseBody []byte

	defer func() {
		log.Context(ctx).
			With("url", url).
			With("method", method).
			With("requestBody", string(requestBody), "responseBody", string(responseBody)).
			Debug()
	}()

	httpClient := xray.Client(&http.Client{Timeout: time.Second * 5})

	request, err := http.NewRequestWithContext(ctx,
		method,
		url,
		bytes.NewBuffer(requestBody),
	)

	if nil != err {
		return 0, nil, err
	}

	request.Header.Set("Content-Type", "application/json")

	for k, v := range headers {
		request.Header.Set(k, v)
	}

	response, err := httpClient.Do(request)
	if err != nil {
		return 0, nil, err
	}

	defer func(body io.ReadCloser) {
		closeErr := body.Close()
		if closeErr != nil {
			log.Context(ctx).Error(closeErr)
		}
	}(response.Body)

	responseBody, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return 0, nil, err
	}

	return response.StatusCode, responseBody, nil
}
