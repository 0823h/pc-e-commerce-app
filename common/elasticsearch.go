package common

import (
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

var ES *elasticsearch.Client

func ESInit() *elasticsearch.Client {
	// cfg := elasticsearch.Config{
	// 	// Addresses: []string{
	// 	// 	"http://localhost:9200",
	// 	// },
	// 	// Transport: &http.Transport{
	// 	// 	MaxIdleConnsPerHost:   10,
	// 	// 	ResponseHeaderTimeout: time.Second,
	// 	// 	TLSClientConfig: &tls.Config{
	// 	// 		MinVersion: tls.VersionTLS12,
	// 	// 		// ...
	// 	// 	},
	// 	// 	// ...
	// 	// },
	// }
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	ES = es
	return ES
}

func GetES() *elasticsearch.Client {
	return ES
}
