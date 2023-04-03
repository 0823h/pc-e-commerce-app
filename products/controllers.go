package products

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"tmdt-backend/common"
	"tmdt-backend/users"

	"github.com/gin-gonic/gin"
)

func GetAllProducts(c *gin.Context) {
	conditionQuery := BindConditionQuery(c)

	var products []Product
	pagination := common.NewPagination()
	common.GetPaginationParameter(c, &pagination)

	db := common.GetDB()
	// db.Where("products.is_deleted = ?", "false").Scopes(common.Paginate(products, &pagination, db)).Joins("Manufacturer").Find(&products)
	db.Where(conditionQuery).Scopes(common.Paginate(products, &pagination, db)).Joins("Manufacturer").Find(&products)
	serializer := ProductsSerializer{c, products}
	pagination.Data = serializer.Response()
	common.SendResponse(c, http.StatusOK, "Success", pagination)
	return
}

// Create new product in database
func CreateProduct(c *gin.Context) {
	validator := NewCreateProductValidator()

	if err := validator.Bind(c); err != nil {
		common.SendResponse(c, http.StatusUnprocessableEntity, err.Error(), nil)
		return
	}

	db := common.GetDB()
	err := db.Create(&validator.productModel).Error
	if err != nil {
		common.SendResponse(c, http.StatusUnprocessableEntity, err.Error(), nil)
		return
	}

	// Create category product relation
	for _, category_id := range validator.CategoriesID {
		var category_product CategoryProductRelation
		category_product.ProductID = uint(validator.productModel.ID)
		category_product.CategoryID = category_id
		if err := db.Create(&category_product).Error; err != nil {
			fmt.Println(err.Error())
		}
	}

	createdProduct := NewProduct()
	db.Where("products.id = ?", validator.productModel.ID).Joins("Manufacturer").First(&createdProduct)

	// Gorse
	// gorse := common.GetGorse()
	// gorse.InsertItem(context.Background(), client.Item{ItemId: strconv.FormatUint(createdProduct.ID, 10), IsHidden: false,
	// 	Categories: []string{"Shoes"}, Timestamp: "2023/03/18 12:22", Labels: []string{"Shoes labels"}, Comment: ""})

	serializer := ProductSerializer{c, createdProduct}
	common.SendResponse(c, http.StatusCreated, "Success", serializer.Response())
	return
}

func UpdateProduct(c *gin.Context) {
	productId := c.Param("id")

	if productId == ":id" {
		common.SendResponse(c, http.StatusBadRequest, "Error: Id not found!", "")
		return
	}

	product := NewProduct()
	db := common.GetDB()
	result := db.Where("products.id = ?", productId).Joins("Manufacturer").First(&product)

	if result.RowsAffected == 0 {
		common.SendResponse(c, http.StatusNotFound, "Error: "+"Product not found", nil)
		return
	}

	updateProductValidator := NewUpdateProductValidatorFillWith(product)
	err := updateProductValidator.Bind(c)
	if err != nil {
		common.SendResponse(c, http.StatusUnprocessableEntity, "Error: "+err.Error(), nil)
		return
	}

	product = updateProductValidator.productModel
	fmt.Println(updateProductValidator)

	db.Save(&product)

	common.SendResponse(c, http.StatusOK, "Success", product)

}

type ConditionQuery struct {
	Name           string
	ID             string
	ManufacturerID string
	IsDeleted      bool
}

func BindConditionQuery(c *gin.Context) ConditionQuery {
	conditionQuery := ConditionQuery{Name: c.Query("name"), ID: c.Query("id"), ManufacturerID: c.Query("manufacturer_id"), IsDeleted: false}
	return conditionQuery
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
	query := QueryBind(c)
	// query := GetESQuery(c)

	var buf bytes.Buffer
	var r map[string]interface{}
	// fmt.Println("QUERY: ", query.SearchKeyword.Name)
	es_query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": query.SearchKeyword,
		},
	}

	// es_query := map[string]interface{}{
	// 	"query": map[string]interface{}{
	// 		"match": query["query"].([]interface{}),
	// 	},
	// }
	if err := json.NewEncoder(&buf).Encode(es_query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("product"),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	// res1 := es.Search().Raw([]byte(`{
	// 	"query": {
	// 	  "term": {
	// 		"user.id": {
	// 		  "value": "kimchy",
	// 		  "boost": 1.0
	// 		}
	// 	  }
	// 	}
	//   }`))

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

func RateProduct(c *gin.Context) {
	// Check product id
	productId := c.Param("id")
	// fmt.Println("Product ID: ", productId)
	if productId == ":id" {
		common.SendResponse(c, http.StatusBadRequest, "Error: Id not found!", "")
		return
	}
	db := common.GetDB()

	product := NewProduct()
	err := db.First(&product, productId).Error
	if err != nil {
		common.SendResponse(c, http.StatusNotFound, "Error: Product not found!", nil)
		return
	}

	// Check user id
	userId := c.GetString("id")
	user := users.NewUser()

	err = db.First(&user, userId).Error
	if err != nil {
		common.SendResponse(c, http.StatusNotFound, "Error: User not found", nil)
		return
	}

	// Bind request
	ratingValidator := NewRatingValidator()
	if err := ratingValidator.Bind(c); err != nil {
		common.SendResponse(c, http.StatusUnprocessableEntity, err.Error(), nil)
		return
	}

	// Search if rating exists
	rating := NewRating()
	result := db.Where("user_id = ? AND product_id = ?", userId, productId).First(&rating)
	if result.RowsAffected == 0 {
		// fmt.Println("Go in")

		userIdUint, _ := strconv.ParseUint(productId, 10, 64)
		productIdUint, _ := strconv.ParseUint(userId, 10, 64)
		rating = Rating{
			Rate:      ratingValidator.Rate,
			ProductID: productIdUint,
			Product:   product,
			UserID:    userIdUint,
			User:      user,
		}
		db.Create(&rating)
		common.SendResponse(c, http.StatusOK, "Success", rating)
		return
	}
	// fmt.Println("Go out", result.RowsAffected)

	result.Update("rate", ratingValidator.Rate)
	common.SendResponse(c, http.StatusOK, "Success", rating)
	return
}

func GetRatings(c *gin.Context) {
	var ratings []Rating
	pagination := common.NewPagination()
	common.GetPaginationParameter(c, &pagination)

	db := common.GetDB()
	db.Where("ratings.is_deleted = ?", "false").Scopes(common.Paginate(ratings, &pagination, db)).Find(&ratings)
	serializer := RatingsSerializer{c, ratings}
	pagination.Data = serializer.Response()
	common.SendResponse(c, http.StatusOK, "Success", pagination)
	return
}

func TestImageUpload(c *gin.Context) {
	file, header, err := c.Request.FormFile("upload")
	filename := header.Filename
	fmt.Println(header.Filename)
	image_link := "./product_images/" + filename + ".png"
	out, err := os.Create(filename)
	// out, err := os.Create("./product_images/" + filename + ".png")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}

	common.SendResponse(c, http.StatusOK, "Success", image_link)
	return
}

func GetCategoryProduct(c *gin.Context) {
	category_id := c.Param("id")
	db := common.GetDB()

	var category_products []CategoryProductRelation
	db.Where("category_product_relations.category_id = ?", category_id).Find(&category_products)

	var product_ids []uint
	for _, category_product := range category_products {
		product_ids = append(product_ids, category_product.ProductID)
	}
	// fmt.Printf("product_ids: %v\n", product_ids)

	var products []Product
	for _, product_id := range product_ids {
		var product Product
		db.Where("products.id = ?", product_id).Joins("Manufacturer").Find(&product)
		products = append(products, product)
		// fmt.Printf("product: %v\n", product)
	}

	common.SendResponse(c, http.StatusOK, "Success", products)
	return
}
