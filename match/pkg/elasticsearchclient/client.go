package elasticsearchclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type Config struct {
	Addresses []string
	Username  string
	Password  string
	// Optional Timeout or other settings can be added
}

type ElasticsearchClient struct {
	serverConf *Config
	client     *elasticsearch.Client
	ctx        context.Context
	cancel     context.CancelFunc
}

func New(conf *Config, opts ...func(*ElasticsearchClient)) *ElasticsearchClient {
	cli := &ElasticsearchClient{serverConf: conf}

	for _, opt := range opts {
		opt(cli)
	}

	return cli
}

func (cli *ElasticsearchClient) Start() error {
	cfg := elasticsearch.Config{
		Addresses: cli.serverConf.Addresses,
	}

	if cli.serverConf.Username != "" {
		cfg.Username = cli.serverConf.Username
		cfg.Password = cli.serverConf.Password
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("elasticsearch client init error: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	res, err := client.Info()
	if err != nil {
		cancel()
		return fmt.Errorf("elasticsearch info check error: %w", err)
	}
	defer res.Body.Close()

	cli.ctx = ctx
	cli.cancel = cancel
	cli.client = client

	return nil
}

func (cli *ElasticsearchClient) Close() {
	if cli.cancel != nil {
		cli.cancel()
	}
}

// Example helper: Index a document
func (cli *ElasticsearchClient) Index(index string, id string, document any) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(document); err != nil {
		return fmt.Errorf("encode error: %w", err)
	}

	res, err := cli.client.Index(
		index,
		&buf,
		cli.client.Index.WithDocumentID(id),
		cli.client.Index.WithContext(cli.ctx),
		cli.client.Index.WithRefresh("true"),
	)
	if err != nil {
		return fmt.Errorf("index error: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("index response error: %s", res.String())
	}

	return nil
}

// Example helper: Search
func (cli *ElasticsearchClient) Search(index string, query map[string]any) (*esapi.Response, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("encode query error: %w", err)
	}

	res, err := cli.client.Search(
		cli.client.Search.WithContext(cli.ctx),
		cli.client.Search.WithIndex(index),
		cli.client.Search.WithBody(&buf),
		cli.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, fmt.Errorf("search error: %w", err)
	}

	if res.IsError() {
		defer res.Body.Close()
		return nil, fmt.Errorf("search response error: %s", res.String())
	}

	return res, nil
}
