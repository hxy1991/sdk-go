package opensearch

import (
	"context"
	"github.com/opensearch-project/opensearch-go"
	"os"
	"testing"
	"time"
)

var client = _defaultClient()

func _defaultClient() *OpenSearch {
	// Initialize the client with SSL/TLS enabled.
	guildIndexName := os.Getenv("GUILD_INDEX_NAME")
	if guildIndexName == "" {
		panic("guildIndexName can not be empty, please set GUILD_INDEX_NAME env")
	}

	_client, err := New(guildIndexName)
	if err != nil {
		panic(err)
	}

	return _client
}

func TestOpenSearch_Scroll(t *testing.T) {
	type fields struct {
		client    *opensearch.Client
		indexName string
	}
	type args struct {
		ctx      context.Context
		jsonBody string
		scroll   time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *SearchResponse
		wantErr bool
	}{
		{
			name: "",
			args: args{
				ctx:      context.TODO(),
				jsonBody: `{"size":20}`,
				scroll:   time.Minute * 3,
			},
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			resp, err := client.SearchWithScroll(tt.args.ctx, tt.args.jsonBody, tt.args.scroll)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scroll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for _, hit := range resp.Hits.Hits {
				t.Log(hit.Source)
			}

			scrollId := resp.ScrollId

			for {
				got, err := client.Scroll(tt.args.ctx, scrollId, tt.args.scroll)
				if (err != nil) != tt.wantErr {
					t.Errorf("Scroll() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				for _, hit := range got.Hits.Hits {
					t.Log(hit.Source)
				}

				scrollId = got.ScrollId

				if len(got.Hits.Hits) <= 0 {
					break
				}
			}

		})
	}
}
