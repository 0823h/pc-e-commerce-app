package products

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"tmdt-backend/common"

	"github.com/gin-gonic/gin"
)

func GetAllProducts(c *gin.Context) {
	var products []Product
	pagination := common.NewPagination()
	common.GetPaginationParameter(c, &pagination)

	db := common.GetDB()
	db.Where("products.is_deleted = ?", "false").Scopes(common.Paginate(products, &pagination, db)).Joins("Manufacturer").Find(&products)
	serializer := ProductsSerializer{c, products}
	pagination.Data = serializer.Response()
	// c.JSON(http.StatusOK, gin.H{"products": products})
	common.SendResponse(c, http.StatusOK, "Success", pagination)
	return
}

func CreateProduct(c *gin.Context) {
	validator := NewCreateProductValidator()

	if err := validator.Bind(c); err != nil {
		// c.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
		// c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		common.SendResponse(c, http.StatusUnprocessableEntity, err.Error(), nil)
		return
	}

	if err := SaveOne(&validator.productModel); err != nil {
		// c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		// c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		common.SendResponse(c, http.StatusUnprocessableEntity, err.Error(), nil)
		return
	}

	serializer := ProductSerializer{c, validator.productModel}
	// c.JSON(http.StatusCreated, gin.H{"Product": serializer.Response()})
	common.SendResponse(c, http.StatusCreated, "Success", serializer.Response())
	return
}

func UpdateProduct(c *gin.Context) {

}

type Query struct {
	Metadata struct {
		page    string
		limit   string
		order   string
		sort    string
		keyword string
	}
	SearchKeyword struct {
		Name string
		// ID   string
		// SKU  string
	}
}

func QueryBind(c *gin.Context) Query {
	page := c.Query("page")
	if page == "" {
		page = "0"
	}
	//p.Metadata.Page, _ = strconv.Atoi(page)

	limit := c.Query("limit")
	if limit == "" {
		limit = "0"
	}
	//p.Metadata.Limit, _ = strconv.Atoi(limit)

	order := strings.ToLower(c.Query("order"))
	if order == "" {
		order = "desc"
	}

	sort := strings.ToLower(c.Query("sort"))
	if sort == "" {
		sort = "id"
	}

	name := strings.ToLower(c.Query("Name"))
	//fmt.Println(c.Query("Name"))

	var query Query
	query.Metadata.page = page
	query.Metadata.limit = limit
	query.Metadata.order = order
	query.Metadata.sort = sort
	query.SearchKeyword.Name = name

	return query
}

func GetESQuery(c *gin.Context) map[string]interface{} {
	var query_array []string
	es_query := make(map[string]interface{})
	name := c.Query("Name")
	if name != "" {
		query_array = append(query_array, "Name")
	}
	id := c.Query("ID")
	if id != "" {
		query_array = append(query_array, "ID")
	}
	sku := c.Query("SKU")
	if sku != "" {
		query_array = append(query_array, "SKU")
	}
	for i := 0; i < len(query_array); i++ {
		es_query["query"].(map[string]interface{})[query_array[i]] = c.Query(strings.ToLower(query_array[i]))
	}

	fmt.Println(es_query)
	return es_query
}

func GetAllProductsES(c *gin.Context) {
	es := common.GetES()
	// query := QueryBind(c)
	query := GetESQuery(c)

	var buf bytes.Buffer
	var r map[string]interface{}
	// fmt.Println("QUERY: ", query.SearchKeyword.Name)
	// es_query := map[string]interface{}{
	// 	"query": map[string]interface{}{
	// 		"match": query.SearchKeyword,
	// 	},
	// }
	es_query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": query["query"].([]interface{}),
		},
	}
	if err := json.NewEncoder(&buf).Encode(es_query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}
	// res, err := es.Search(
	// 	es.Search.WithContext(context.Background()),
	// 	es.Search.WithIndex("product"),
	// 	es.Search.WithBody(&buf),
	// 	es.Search.WithTrackTotalHits(true),
	// 	es.Search.WithPretty(),
	// )
	// if err != nil {
	// 	log.Fatalf("Error getting response: %s", err)
	// }

	res1 := es.Search().Raw([]byte(`{
		"query": {
		  "term": {
			"user.id": {
			  "value": "kimchy",
			  "boost": 1.0
			}
		  }
		}
	  }`))

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Printf("Error parsing the response body: %s", err)
		common.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Printf("Error parsing the response body: %s", err)
			common.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
		} else {
			// Print the response status and error information.
			log.Fatalf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	log.Printf(
		"[%s] %d hits; took: %dms",
		res.Status(),
		int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		int(r["took"].(float64)),
	)

	fmt.Println(r["hits"].(map[string]interface{})["hits"])
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		log.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
	}

	log.Println(strings.Repeat("=", 37))

	// common.SendResponse(c, http.StatusOK, "Success", nil)
	common.SendResponse(c, http.StatusOK, "Success", r["hits"].(map[string]interface{})["hits"].([]interface{}))
	return
}
