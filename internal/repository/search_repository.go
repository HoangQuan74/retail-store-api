package repository

import (
	"github.com/elastic/go-elasticsearch/v8"
	es "github.com/kainguyen/retail-store-api/pkg/elasticsearch"
)

type SearchRepository struct {
	client    *elasticsearch.Client
	indexName string
}

func NewSearchRepository(client *elasticsearch.Client, indexName string) *SearchRepository {
	return &SearchRepository{client: client, indexName: indexName}
}

func (r *SearchRepository) IndexProduct(doc es.ProductDocument) error {
	return es.IndexProduct(r.client, r.indexName, doc)
}

func (r *SearchRepository) DeleteProduct(id int64) error {
	return es.DeleteProduct(r.client, r.indexName, id)
}

func (r *SearchRepository) Search(query string, limit, offset int) (*es.SearchResult, error) {
	return es.SearchProducts(r.client, r.indexName, query, limit, offset)
}
