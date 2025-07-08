// Package config config/elasticsearch.go
package config

import (
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"net/http"
	"time"
)

var ESClient *elasticsearch.Client

// InitElasticsearch initializes and returns the Elasticsearch client
func InitElasticsearch() {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
		Transport: &http.Transport{
			ResponseHeaderTimeout: 2 * time.Second,
		},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the Elasticsearch client: %s", err)
	}

	ESClient = es
	log.Println("ES connected successfully")
}
