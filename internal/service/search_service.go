package service

import (
	"github.com/elastic/go-elasticsearch/v8"
	es "github.com/kainguyen/retail-store-api/pkg/elasticsearch"
	"github.com/kainguyen/retail-store-api/internal/repository"
)

type SearchService struct {
	repo *repository.SearchRepository
}

func NewSearchService(client *elasticsearch.Client, indexName string) *SearchService {
	return &SearchService{repo: repository.NewSearchRepository(client, indexName)}
}

func (s *SearchService) SearchProducts(query string, limit, offset int) (*es.SearchResult, error) {
	return s.repo.Search(query, limit, offset)
}
