package http

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/hxy1991/sdk-go/log"
	"github.com/hxy1991/sdk-go/utils"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

// Defaults for the HTTPTransportBuilder.
var (
	DefaultTimeout = 5 * time.Second

	// Default connection pool options
	DefaultHTTPTransportMaxIdleConns        = 200
	DefaultHTTPTransportMaxIdleConnsPerHost = 50

	// Default connection timeouts
	DefaultHTTPTransportIdleConnTimeout       = 90 * time.Second
	DefaultHTTPTransportTLSHandleshakeTimeout = 10 * time.Second
	DefaultHTTPTransportExpectContinueTimeout = 1 * time.Second

	// Default to TLS 1.2 for all HTTPS requests.
	DefaultHTTPTransportTLSMinVersion uint16 = tls.VersionTLS12
)

// Timeouts for net.Dialer's network connection.
var (
	DefaultDialConnectTimeout   = 30 * time.Second
	DefaultDialKeepAliveTimeout = 30 * time.Second
)

func defaultHTTPTransport() *http.Transport {
	dialer := &net.Dialer{
		Timeout:   DefaultDialConnectTimeout,
		KeepAlive: DefaultDialKeepAliveTimeout,
		DualStack: true,
	}

	tr := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		TLSHandshakeTimeout:   DefaultHTTPTransportTLSHandleshakeTimeout,
		MaxIdleConns:          DefaultHTTPTransportMaxIdleConns,
		MaxIdleConnsPerHost:   DefaultHTTPTransportMaxIdleConnsPerHost,
		IdleConnTimeout:       DefaultHTTPTransportIdleConnTimeout,
		ExpectContinueTimeout: DefaultHTTPTransportExpectContinueTimeout,
		ForceAttemptHTTP2:     true,
		TLSClientConfig: &tls.Config{
			MinVersion: DefaultHTTPTransportTLSMinVersion,
		},
	}

	return tr
}

var client = xray.Client(&http.Client{
	Timeout:   DefaultTimeout,
	Transport: defaultHTTPTransport(),
})

func Send(ctx context.Context, url, method string, requestBody []byte, headers map[string]string) (int, []byte, error) {
	return SendWithTimeout(ctx, url, method, requestBody, headers, 5, 10*time.Millisecond)
}

func SendWithLogMinDuration(ctx context.Context, url, method string, requestBody []byte, headers map[string]string, logMinDuration time.Duration) (int, []byte, error) {
	return SendWithTimeout(ctx, url, method, requestBody, headers, 5, logMinDuration)
}

func SendWithTimeout(ctx context.Context, url, method string, requestBody []byte, headers map[string]string, second int,
	logMinDuration time.Duration) (int, []byte, error) {
	startTime := time.Now()

	var responseBody []byte
	var responseCode int

	defer func() {
		duration := time.Now().Sub(startTime)
		if duration > logMinDuration || responseCode != 200 || !utils.IsProduction() {
			log.Context(ctx).
				With("requestPath", url).
				With("requestMethod", method).
				With("responseCode", responseCode).
				With("requestBody", string(requestBody), "responseBody", string(responseBody)).
				With("requestHeader", headers).
				With("latency", fmt.Sprintf("%13v", duration)).
				With("latencyInNS", duration.Nanoseconds()).
				Debug()
		}
	}()

	newCtx, cancel := context.WithTimeout(ctx, time.Duration(second)*time.Second)
	defer cancel()

	request, err := http.NewRequestWithContext(newCtx,
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

	response, err := client.Do(request)
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

	responseCode = response.StatusCode

	return response.StatusCode, responseBody, nil
}
