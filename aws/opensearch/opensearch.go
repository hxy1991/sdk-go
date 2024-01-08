package opensearch

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/hxy1991/sdk-go/aws/secretsmanager"
	"github.com/hxy1991/sdk-go/log"
	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
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

func (openSearch *OpenSearch) GetClient() *opensearch.Client {
	return openSearch.client
}

func (openSearch *OpenSearch) GetIndexName() string {
	return openSearch.indexName
}

type config struct {
	domainEndpoint string
	username       string
	password       string
}

func New(ctx context.Context, indexName string) (openSearch *OpenSearch, err error) {
	openSearchConfig := getOpenSearchConfig(ctx)

	var transport http.RoundTripper
	transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// https://opensearch.org/docs/latest/clients/go/
	client, err := opensearch.NewClient(opensearch.Config{
		Transport:            xray.RoundTripper(transport),
		Addresses:            []string{openSearchConfig.domainEndpoint},
		Username:             openSearchConfig.username,
		Password:             openSearchConfig.password,
		UseResponseCheckOnly: true,
	})
	if err != nil {
		return nil, err
	}
	return &OpenSearch{client: client, indexName: indexName}, nil
}

func NewWithAuth(ctx context.Context, indexName, domainEndpoint, username, password string) (openSearch *OpenSearch, err error) {
	var transport http.RoundTripper
	transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// https://opensearch.org/docs/latest/clients/go/
	client, err := opensearch.NewClient(opensearch.Config{
		Transport:            xray.RoundTripper(transport),
		Addresses:            []string{domainEndpoint},
		Username:             username,
		Password:             password,
		UseResponseCheckOnly: true,
	})
	if err != nil {
		return nil, err
	}
	return &OpenSearch{client: client, indexName: indexName}, nil
}

func NewWithEndpoint(ctx context.Context, domainEndpoint, indexName, username, password string) (openSearch *OpenSearch, err error) {
	var transport http.RoundTripper
	transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// https://opensearch.org/docs/latest/clients/go/
	client, err := opensearch.NewClient(opensearch.Config{
		Transport:            xray.RoundTripper(transport),
		Addresses:            []string{domainEndpoint},
		Username:             username,
		Password:             password,
		UseResponseCheckOnly: true,
	})
	if err != nil {
		return nil, err
	}
	return &OpenSearch{client: client, indexName: indexName}, nil
}

func getOpenSearchConfig(ctx context.Context) *config {
	domainEndpoint := os.Getenv("OPEN_SEARCH_DOMAIN_ENDPOINT")
	if domainEndpoint == "" {
		panic("openSearch domainEndpoint can not be empty, please set OPEN_SEARCH_DOMAIN_ENDPOINT env")
	} else if !strings.HasPrefix(strings.ToLower(domainEndpoint), "http") {
		domainEndpoint = "https://" + domainEndpoint
	}

	username := os.Getenv("OPEN_SEARCH_MASTER_USER_NAME")
	if username == "" {
		panic("openSearch master username can not be empty, please set OPEN_SEARCH_MASTER_USER_NAME env")
	}

	passwordSecretName := os.Getenv("OPEN_SEARCH_MASTER_PASSWORD_SECRET_NAME")
	if passwordSecretName == "" {
		panic("openSearch passwordSecretName can not be empty, please set OPEN_SEARCH_MASTER_PASSWORD_SECRET_NAME env")
	}

	return &config{
		domainEndpoint: domainEndpoint,
		username:       username,
		password:       *getPassword(ctx, passwordSecretName),
	}
}

func getPassword(ctx context.Context, passwordSecretName string) *string {
	getSecretValueOutput, err := awssecretmanager.GetSecret(ctx, passwordSecretName)
	if err != nil {
		panic(err)
	}
	return getSecretValueOutput.SecretString
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

	body, err := ioutil.ReadAll(response.Body)
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

	body, err := ioutil.ReadAll(response.Body)
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

	body, err := ioutil.ReadAll(response.Body)
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
