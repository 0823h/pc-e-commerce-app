package elasticsearch

import (
	"log"
	"strings"
	"tmdt-backend/common"

	"github.com/gin-gonic/gin"
)

func ElasticRouter(router *gin.RouterGroup) {
	router.GET("/search")
}

func CreateIndex(c *gin.Context) {
	es := common.ESInit()

	index := "products"
	mapping := `{
		"settings": {
			"number_of_shards": 1
		},
		"mappings": {
			"properties": {
				"field1": {
					"type": "text"
				}
			}
		}
	}`

	res, err := es.Indices.Create(
		index,
		es.Indices.Create.WithBody(strings.NewReader(mapping)),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(res)
}
