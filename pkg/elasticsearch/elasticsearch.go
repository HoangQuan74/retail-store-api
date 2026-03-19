package elasticsearch

import (
	"fmt"
	"log/slog"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/kainguyen/retail-store-api/internal/config"
)

func NewClient(cfg config.ElasticsearchConfig) (*elasticsearch.Client, error) {
	esCfg := elasticsearch.Config{
		Addresses: cfg.Addresses,
	}
	if cfg.Username != "" {
		esCfg.Username = cfg.Username
		esCfg.Password = cfg.Password
	}

	client, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		return nil, fmt.Errorf("create elasticsearch client: %w", err)
	}

	res, err := client.Ping()
	if err != nil {
		return nil, fmt.Errorf("ping elasticsearch: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch ping error: %s", res.String())
	}

	slog.Info("Connected to Elasticsearch successfully")
	return client, nil
}
