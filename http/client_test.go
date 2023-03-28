package http

import (
	"context"
	"fmt"
	"testing"
)

func TestSend(t *testing.T) {
	type args struct {
		ctx         context.Context
		url         string
		method      string
		requestBody []byte
		headers     map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "",
			args: args{
				ctx:         context.TODO(),
				url:         "https://www.google.com",
				method:      "GET",
				requestBody: nil,
				headers:     nil,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode, got, err := Send(tt.args.ctx, tt.args.url, tt.args.method, tt.args.requestBody, tt.args.headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			fmt.Println(statusCode, " ", string(got))
		})
	}
}

func TestSendWithTimeout(t *testing.T) {
	type args struct {
		ctx         context.Context
		url         string
		method      string
		requestBody []byte
		headers     map[string]string
		second      int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		want1   []byte
		wantErr bool
	}{
		{
			name: "",
			args: args{
				ctx:         context.TODO(),
				url:         "https://www.google.com",
				method:      "GET",
				requestBody: nil,
				headers:     nil,
				second:      3,
			},
			want:    200,
			want1:   nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, err := SendWithTimeout(tt.args.ctx, tt.args.url, tt.args.method, tt.args.requestBody, tt.args.headers, tt.args.second)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendWithTimeout() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SendWithTimeout() got = %v, want %v", got, tt.want)
			}
		})
	}
}
