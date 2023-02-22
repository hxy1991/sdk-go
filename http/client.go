package http

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/hxy1991/sdk-go/log"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

// Defaults for the HTTPTransportBuilder.
var (
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

func defaultDialer() *net.Dialer {
	return &net.Dialer{
		Timeout:   DefaultDialConnectTimeout,
		KeepAlive: DefaultDialKeepAliveTimeout,
		DualStack: true,
	}
}

func defaultHTTPTransport() *http.Transport {
	dialer := defaultDialer()

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

func Send(ctx context.Context, url, method string, requestBody []byte, headers map[string]string) (int, []byte, error) {
	return SendWithTimeout(ctx, url, method, requestBody, headers, 5)
}

func SendWithTimeout(ctx context.Context, url, method string, requestBody []byte, headers map[string]string, second int) (int, []byte, error) {
	startTime := time.Now()

	var responseBody []byte

	defer func() {
		log.Context(ctx).
			With("url", url).
			With("method", method).
			With("requestBody", string(requestBody), "responseBody", string(responseBody)).
			With("latency", fmt.Sprintf("%13v", time.Now().Sub(startTime))).
			Debug()
	}()

	client := &http.Client{Timeout: time.Duration(second) * time.Second}
	client.Transport = defaultHTTPTransport()

	httpClient := xray.Client(client)

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
