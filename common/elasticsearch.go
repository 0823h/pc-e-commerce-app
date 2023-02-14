package common

import (
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

var ES *elasticsearch.Client

func ESInit() *elasticsearch.Client {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// res, err := es.Info()
	// if err != nil {
	// 	log.Println(err.Error())
	// }

	// defer res.Body.Close()
	// log.Println(res)

	ES = es

	return ES
}

func GetES() *elasticsearch.Client {
	return ES
}
