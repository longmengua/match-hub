package elasticsearchclient_test

import (
	"encoding/json"
	"io"
	"match/pkg/elasticsearchclient"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestElasticsearchClient_Integration(t *testing.T) {
	conf := &elasticsearchclient.Config{
		Addresses: []string{"http://elasticsearch.sql.orb.local:9200"},
	}

	client := elasticsearchclient.New(conf)
	err := client.Start()
	assert.NoError(t, err, "Start should not return error")
	defer client.Close()

	// Index a document
	docID := "test-doc-1"
	err = client.Index("test-index", docID, map[string]interface{}{
		"title": "Integration Test",
		"ts":    time.Now(),
	})
	assert.NoError(t, err, "Index should not return error")

	// Search for the document
	res, err := client.Search("test-index", map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"title": "Integration",
			},
		},
	})
	assert.NoError(t, err, "Search should not return error")
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	assert.NoError(t, err, "Should read search response body")

	var searchResult map[string]interface{}
	err = json.Unmarshal(bodyBytes, &searchResult)
	assert.NoError(t, err, "Should unmarshal search result")

	hits := searchResult["hits"].(map[string]interface{})["hits"].([]interface{})
	assert.True(t, len(hits) >= 1, "Should have at least one search hit")
}
