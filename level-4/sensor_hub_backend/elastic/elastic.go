package elastic

import (
	"fmt"
	"os"

	"github.com/elastic/go-elasticsearch/v9"
)

var es *elasticsearch.Client

func InitConnection() {
	url := os.Getenv("ELASTIC_URL")
	key := os.Getenv("ELASTIC_API_KEY")

	esConfig := elasticsearch.Config{
		Addresses: []string{url},
		APIKey:    key,
		//EnableMetrics: true, // TODO enable?
	}

	client, err := elasticsearch.NewClient(esConfig)
	if err != nil {
		panic(err)
	}

	es = client

	fmt.Println("Connected to ElasticSearch")
}
