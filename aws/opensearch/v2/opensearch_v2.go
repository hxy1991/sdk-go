package opensearchv2

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/hxy1991/sdk-go/log"
	"github.com/opensearch-project/opensearch-go/v2"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
	"github.com/opensearch-project/opensearch-go/v2/signer/awsv2"
)

type SearchResponse struct {
	Took     int     `json:"took"`
	TimedOut bool    `json:"timed_out"`
	Shards   *Shards `json:"_shards"`
	Hits     *Hits   `json:"hits"`
	ScrollId string  `json:"_scroll_id"`
}
type Shards struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Skipped    int `json:"skipped"`
	Failed     int `json:"failed"`
}
type Total struct {
	Value    int64  `json:"value"`
	Relation string `json:"relation"`
}

type Hits struct {
	Total    *Total       `json:"total"`
	MaxScore float64      `json:"max_score"`
	Hits     []*InnerHits `json:"hits"`
}
type InnerHits struct {
	Index  string      `json:"_index"`
	Type   string      `json:"_type"`
	ID     string      `json:"_id"`
	Score  float64     `json:"_score"`
	Source interface{} `json:"_source"`
}

type ClearScrollResponse struct {
	Succeeded bool `json:"succeeded"`
	NumFreed  int  `json:"num_freed"`
}

type OpenSearch struct {
	client    *opensearch.Client
	indexName string
}

func New(ctx context.Context, indexName string) (openSearch *OpenSearch, err error) {
	region := os.Getenv("AWS_REGION")
	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)
	if err != nil {
		return nil, err
	}

	// Create an AWS request Signer and load AWS configuration using default config folder or env vars.
	signer, err := awsv2.NewSignerWithService(awsCfg, "aoss") // "aoss" for Amazon OpenSearch Serverless
	if err != nil {
		return nil, err
	}

	domainEndpoint := os.Getenv("OPEN_SEARCH_DOMAIN_ENDPOINT")
	if domainEndpoint == "" {
		panic("openSearch domainEndpoint can not be empty, please set OPEN_SEARCH_DOMAIN_ENDPOINT env")
	} else if !strings.HasPrefix(strings.ToLower(domainEndpoint), "http") {
		domainEndpoint = "https://" + domainEndpoint
	}

	client, err := opensearch.NewClient(opensearch.Config{
		Addresses: []string{domainEndpoint},
		Signer:    signer,
	})

	if err != nil {
		return nil, err
	}

	return &OpenSearch{client: client, indexName: indexName}, nil
}

func (openSearch *OpenSearch) CreateIndex(ctx context.Context, json string) (*opensearchapi.Response, error) {
	mapping := strings.NewReader(json)
	req := opensearchapi.IndicesCreateRequest{
		Index: openSearch.indexName,
		Body:  mapping,
	}
	createIndexResponse, err := req.Do(ctx, openSearch.client)
	if err != nil {
		return nil, err
	}
	return createIndexResponse, nil
}

func (openSearch *OpenSearch) AddDoc(ctx context.Context, docId string, json string) (*opensearchapi.Response, error) {
	document := strings.NewReader(json)
	req := opensearchapi.IndexRequest{
		Index:      openSearch.indexName,
		DocumentID: docId,
		Body:       document,
	}
	insertResponse, err := req.Do(ctx, openSearch.client)
	if err != nil {
		return nil, err
	}
	return insertResponse, nil
}

func (openSearch *OpenSearch) DeleteDoc(ctx context.Context, docId string) (*opensearchapi.Response, error) {
	deleteRequest := opensearchapi.DeleteRequest{
		Index:      openSearch.indexName,
		DocumentID: docId,
	}

	deleteResponse, err := deleteRequest.Do(ctx, openSearch.client)
	if err != nil {
		return nil, err
	}
	return deleteResponse, nil
}

func (openSearch *OpenSearch) DeleteIndex(ctx context.Context) (*opensearchapi.Response, error) {
	indicesDeleteRequest := opensearchapi.IndicesDeleteRequest{
		Index: []string{openSearch.indexName},
	}

	deleteIndexResponse, err := indicesDeleteRequest.Do(ctx, openSearch.client)
	if err != nil {
		return nil, err
	}
	return deleteIndexResponse, nil
}

func (openSearch *OpenSearch) Search(ctx context.Context, searchJSONBody string) (*SearchResponse, error) {
	return openSearch.SearchWithScroll(ctx, searchJSONBody, time.Duration(0))
}

func (openSearch *OpenSearch) SearchWithScroll(ctx context.Context, searchJSONBody string, scroll time.Duration) (*SearchResponse, error) {
	content := strings.NewReader(searchJSONBody)

	searchRequest := opensearchapi.SearchRequest{
		Index:  []string{openSearch.indexName},
		Body:   content,
		Scroll: scroll,
	}

	response, err := searchRequest.Do(ctx, openSearch.client)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(response.String())
	}

	defer func(Body io.ReadCloser) {
		closeErr := Body.Close()
		if closeErr != nil {
			log.Error(closeErr)
		}
	}(response.Body)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var searchResponse SearchResponse
	err = json.Unmarshal(body, &searchResponse)
	if err != nil {
		return nil, err
	}
	return &searchResponse, nil

}

// Scroll scroll 参数告诉 OpenSearch 把搜索上下文再保持多久时间
func (openSearch *OpenSearch) Scroll(ctx context.Context, scrollID string, scroll time.Duration) (*SearchResponse, error) {
	scrollRequest := opensearchapi.ScrollRequest{
		ScrollID: scrollID,
		Scroll:   scroll,
	}

	response, err := scrollRequest.Do(ctx, openSearch.client)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(response.String())
	}

	defer func(Body io.ReadCloser) {
		closeErr := Body.Close()
		if closeErr != nil {
			log.Error(closeErr)
		}
	}(response.Body)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var searchResponse SearchResponse
	err = json.Unmarshal(body, &searchResponse)
	if err != nil {
		return nil, err
	}
	return &searchResponse, nil
}

// ClearScroll scroll 参数告诉 OpenSearch 把搜索上下文再保持多久时间
func (openSearch *OpenSearch) ClearScroll(ctx context.Context, scrollId []string) (*ClearScrollResponse, error) {
	scrollRequest := opensearchapi.ClearScrollRequest{
		ScrollID: scrollId,
	}

	response, err := scrollRequest.Do(ctx, openSearch.client)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(response.String())
	}

	defer func(Body io.ReadCloser) {
		closeErr := Body.Close()
		if closeErr != nil {
			log.Error(closeErr)
		}
	}(response.Body)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var clearScrollResponse ClearScrollResponse
	err = json.Unmarshal(body, &clearScrollResponse)
	if err != nil {
		return nil, err
	}
	return &clearScrollResponse, nil
}
