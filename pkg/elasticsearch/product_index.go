package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
)

type ProductDocument struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description,omitempty"`
	Price        float64   `json:"price"`
	Quantity     int32     `json:"quantity"`
	CategoryID   int64     `json:"category_id,omitempty"`
	CategoryName string    `json:"category_name,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

const productMapping = `{
  "mappings": {
    "properties": {
      "id":            { "type": "long" },
      "name":          { "type": "text", "fields": { "keyword": { "type": "keyword" } } },
      "description":   { "type": "text" },
      "price":         { "type": "float" },
      "quantity":      { "type": "integer" },
      "category_id":   { "type": "long" },
      "category_name": { "type": "keyword" },
      "created_at":    { "type": "date" },
      "updated_at":    { "type": "date" }
    }
  }
}`

func EnsureProductIndex(client *elasticsearch.Client, indexName string) error {
	res, err := client.Indices.Exists([]string{indexName})
	if err != nil {
		return fmt.Errorf("check index exists: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		slog.Info("Elasticsearch index already exists", "index", indexName)
		return nil
	}

	res, err = client.Indices.Create(
		indexName,
		client.Indices.Create.WithBody(strings.NewReader(productMapping)),
	)
	if err != nil {
		return fmt.Errorf("create index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("create index error: %s", res.String())
	}

	slog.Info("Elasticsearch index created", "index", indexName)
	return nil
}

func IndexProduct(client *elasticsearch.Client, indexName string, doc ProductDocument) error {
	data, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("marshal product: %w", err)
	}

	res, err := client.Index(
		indexName,
		strings.NewReader(string(data)),
		client.Index.WithDocumentID(fmt.Sprintf("%d", doc.ID)),
		client.Index.WithContext(context.Background()),
	)
	if err != nil {
		return fmt.Errorf("index product: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("index product error: %s", res.String())
	}
	return nil
}

func DeleteProduct(client *elasticsearch.Client, indexName string, id int64) error {
	res, err := client.Delete(
		indexName,
		fmt.Sprintf("%d", id),
		client.Delete.WithContext(context.Background()),
	)
	if err != nil {
		return fmt.Errorf("delete product: %w", err)
	}
	defer res.Body.Close()
	return nil
}

type SearchResult struct {
	Items []ProductDocument
	Total int64
}

func SearchProducts(client *elasticsearch.Client, indexName string, query string, limit, offset int) (*SearchResult, error) {
	searchBody := map[string]interface{}{
		"from": offset,
		"size": limit,
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":     query,
				"fields":    []string{"name^3", "description", "category_name"},
				"fuzziness": "AUTO",
			},
		},
	}

	data, _ := json.Marshal(searchBody)

	res, err := client.Search(
		client.Search.WithContext(context.Background()),
		client.Search.WithIndex(indexName),
		client.Search.WithBody(strings.NewReader(string(data))),
	)
	if err != nil {
		return nil, fmt.Errorf("search products: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("search error: %s", res.String())
	}

	var result struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source ProductDocument `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode search response: %w", err)
	}

	items := make([]ProductDocument, 0, len(result.Hits.Hits))
	for _, hit := range result.Hits.Hits {
		items = append(items, hit.Source)
	}

	return &SearchResult{
		Items: items,
		Total: result.Hits.Total.Value,
	}, nil
}
