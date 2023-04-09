package recommendation

import (
	"fmt"
	"log"
	"math"
	"sort"
	"strconv"
	"tmdt-backend/categories"
	"tmdt-backend/common"
	"tmdt-backend/products"
	"tmdt-backend/users"

	"github.com/gaspiman/cosine_similarity"
	linearmodel "github.com/pa-m/sklearn/linear_model"
	"gonum.org/v1/gonum/mat"
)

type MatrixSlice struct {
	data [][]float64
}

// Creating a new matrix
func NewMatrix() MatrixSlice {
	// data := [][]float64{}
	// matrix := MatrixSlice{data}
	matrix := MatrixSlice{nil}
	return matrix
}

// Insert a row into matrix
func (self *MatrixSlice) InsertRow(row []float64) {
	self.data = append(self.data, row[:])
}

// Get number of rows of matrix
func (self *MatrixSlice) GetNumberOfRows() int {
	return len(self.data)
}

// Get number of columns of matrix
func (self *MatrixSlice) GetNumberOfColumns() int {
	if self.GetNumberOfRows() == 0 {
		return 0
	}
	return len(self.data[0])
}

// Creating a slice from inserted float64-type elements
func ToSlice(elems ...float64) []float64 {
	var slice []float64
	for _, elem := range elems {
		slice = append(slice, elem)
	}
	return slice
}

// Load ratings from database into matrix
func (self *MatrixSlice) LoadRating() {
	db := common.GetDB()
	var ratings []products.Rating
	rowCount := db.Find(&ratings).RowsAffected

	var row []float64
	for i := int64(0); i < rowCount; i++ {
		row = ToSlice(float64(ratings[i].UserID),
			float64(ratings[i].ProductID),
			float64(ratings[i].Rate))
		self.InsertRow(row)
	}

	fmt.Printf("self.data: %v\n", self.data)
}

// Helper function: Find if an element exists in slice
func FindInSlice(slice []float64, x float64) bool {
	for _, elem := range slice {
		if elem == x {
			return true
		}
	}
	return false
}

// Get number of users from a matrix
func (self *MatrixSlice) GetUsersLength() (int, []float64) {
	var usersLength []float64
	for i := 0; i < self.GetNumberOfRows(); i++ {
		if FindInSlice(usersLength, self.data[i][0]) == false {
			usersLength = append(usersLength, self.data[i][0])
		}
	}
	return len(usersLength), usersLength
}

// Get number of items from a matrix
func (self *MatrixSlice) GetItemsLength() int {
	var itemsLength []float64
	for i := 0; i < self.GetNumberOfRows(); i++ {
		if FindInSlice(itemsLength, self.data[i][1]) == false {
			itemsLength = append(itemsLength, self.data[i][1])
		}
	}
	return len(itemsLength)
}

// Get rows of matrix
func (self *MatrixSlice) GetRow(index int) []float64 {
	if index < 0 || index > self.GetNumberOfRows()-1 {
		log.Fatalf("Row index is out of range")
	}
	return self.data[index]
}

// Get columns of matrix
func (self *MatrixSlice) GetColumn(index int) []float64 {
	if index < 0 || index > self.GetNumberOfColumns()-1 {
		log.Fatalf("Column index is out of range")
	}
	var column []float64
	for _, row := range self.data {
		column = append(column, row[index])
	}
	return column
}

// Calculate mean ratings of each users, store into slice
func (self *MatrixSlice) CalculateMeanFromMatrix() [][2]float64 {
	var slice [][2]float64
	n_users, n_users_slice := self.GetUsersLength()
	for i := 0; i < n_users; i++ {
		var mean float64
		var count int = 0
		for j := 0; j < self.GetNumberOfRows(); j++ {
			if self.GetRow(j)[0] == n_users_slice[i] {
				mean += self.GetRow(j)[2]
				count++
			}
		}
		mean /= float64(count)
		slice = append(slice, [2]float64{n_users_slice[i], mean})
	}
	return slice
}

// Get matrix value
func (self *MatrixSlice) GetMatrixValue(row int, column int) float64 {
	return self.data[row][column]
}

type CF struct {
	Y_data    MatrixSlice
	Ybar_data MatrixSlice
	n_users   int
	n_items   int
	k         int
}

// Creating a new CF
func NewCF() CF {
	y_data := NewMatrix()
	y_data.LoadRating()
	ybar_data_data := make([][]float64, len(y_data.data))

	for i := 0; i < y_data.GetNumberOfRows(); i++ {
		// ybar_data.data = append(ybar_data.data, y_data.data[i])
		// 	// copy(ybar_data.data[i], y_data.data[i])
		ybar_data_data[i] = make([]float64, len(y_data.data[i]))
		copy(ybar_data_data[i], y_data.data[i])
	}

	ybar_data := MatrixSlice{ybar_data_data}

	n_users, _ := y_data.GetUsersLength()
	n_items := y_data.GetItemsLength()
	cf := CF{y_data, ybar_data, n_users, n_items, 2}
	return cf
}

// Normalize a maitrx (collaborative filtering)
func (self *CF) Normalize_Y() {
	mean_users := self.Y_data.CalculateMeanFromMatrix()

	for i := 0; i < self.Ybar_data.GetNumberOfRows(); i++ {
		for j := 0; j < len(mean_users); j++ {
			if self.Ybar_data.data[i][0] == mean_users[j][0] {
				self.Ybar_data.data[i][2] -= mean_users[j][1]
			}
		}

	}
}

// Form matrix Y from Y_data
func (self *CF) FormMatrix() {
	n_users, n_products := GetDBNumberOfUsersAndProducts()
	y_matrix := make([][]float64, n_products)
	for i := range y_matrix {
		y_matrix[i] = make([]float64, n_users)
	}

	for i := 0; i < self.Ybar_data.GetNumberOfRows(); i++ {
		// Need to -1 because gorm auto increment starts at 1 (not 0), so UserID 1 will be stored at column 0,
		// and ProductID 1 will stored at row 0
		y_matrix[int(self.Ybar_data.data[i][1])-1][int(self.Ybar_data.data[i][0])-1] = self.Ybar_data.data[i][2]
	}
	self.Ybar_data.data = y_matrix
}

// A main function for testing in Recommendation package
func InitMatrix() {
	cf := NewCF()
	// mean_users := cf.Y_data.CalculateMeanFromMatrix()
	cf.Normalize_Y()
	cf.FormMatrix()
	// user_vector := cf.Ybar_data.GetUserVector(1)
	cf.Recommend(1)

	category_vector := GetCategoryVectors()
	fmt.Printf("category_vector: %v\n", category_vector)

	matrix := GetItemFeaturesMatrix()
	fmt.Printf("matrix: %v\n", matrix)

	rated_matrix, lib_rated_matrix := GetProductsRatedByUser(1)
	fmt.Printf("rated_matrix: %v\n", rated_matrix)
	fmt.Printf("lib_rated_matrix: %v\n", lib_rated_matrix)

	rated_feature_matrix := GetRatedItemFeaturesMatrix(rated_matrix)
	fmt.Printf("rated_feature_matrix: %v\n", rated_feature_matrix)

	clf := linearmodel.NewRidge()
	clf.Tol = 1e-3
	clf.Normalize = false
	clf.Alpha = 1
	clf.L1Ratio = 0.

	clf.Fit(rated_feature_matrix, lib_rated_matrix)
	fmt.Printf("Coef:\n%.2f\n", mat.Formatted(clf.Coef.T()))
	fmt.Printf("Intercept:\n%.2f\n", mat.Formatted(clf.Intercept.T()))
}

// Helper function: Find number of users, number of products
// from database
func GetDBNumberOfUsersAndProducts() (n_users int, n_products int) {
	db := common.GetDB()

	var users []users.User
	var products []products.Product

	result := db.Find(&users)
	n_users = int(result.RowsAffected)

	result = db.Find(&products)
	n_products = int(result.RowsAffected)

	return n_users, n_products
}

// Similarity function to calculate similarity between user - user using cosine - similarity
func CalculateSimilarity(user_1 []float64, user_2 []float64) float64 {
	ProcessIfZeroVector(user_1)
	ProcessIfZeroVector(user_2)

	distant, err := cosine_similarity.Cosine(user_1, user_2)
	if err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
		log.Panic("Calculate distance fail!")
	}
	return distant
}

// Prediction function
func (self *CF) Predict(user_id int, item_id int) float64 {
	// Get ids of all users who rate the item, has form [[other_user_id rating_on_item], ...]
	users_slice := self.Y_data.FindUsersWhoRateItem(float64(item_id), float64(user_id))
	fmt.Printf("users_slice: %v\n", users_slice)
	// Get vector of user for cosine similarity, has form [[rating_on_item1 rating_on_item2 rating_on_item3 ...]]
	user_vector := self.Ybar_data.GetUserVector(float64(user_id))
	// Slice to store all calculated similarity, has form [[other_user_id similarity], ...]
	var similarity_slice [][]float64
	// Calculate similarity and append into the slice above
	for _, other_user := range users_slice {
		// Get vector of other user, has form [[rating_on_item1 rating_on_item2 rating_on_item3 ...]]
		other_user_vector := self.Ybar_data.GetUserVector(other_user[0])
		fmt.Printf("other_user_vector: %v\n", other_user_vector)
		// Creating a temporary slice for storing one calculated similarity above, has form [other_user_id similarity]
		temp_slice := []float64{other_user[0], CalculateSimilarity(user_vector, other_user_vector)}
		// Append into similarity_slice
		similarity_slice = append(similarity_slice, temp_slice)
	}

	fmt.Printf("similarity_slice before select: %v\n", similarity_slice)
	// Sort similarity_slice following decrease order of similarity
	sort.Slice(similarity_slice, func(i, j int) bool {
		return similarity_slice[i][1] > similarity_slice[j][1]
	})
	// Select first n-th elements which is first n-th other-user that is most similar to user (n is k specified in struct CF)
	similarity_slice = SelectFirstElementsOfSlice(similarity_slice, self.k)
	fmt.Printf("similarity_slice: %v\n", similarity_slice)
	// Get normalized rating of them on item
	// Create a slice to store normalize_rating, has form [[other_user_id normalized_rating], ...]
	var normalized_rating_slice [][]float64
	// Loop through similarity slice to get other user id and get normalize rating
	for _, other_user := range similarity_slice {
		// Get normalize rating of other user, has form [other_user_id normalized_rating]
		normalize_rating := []float64{other_user[1],
			self.Ybar_data.GetMatrixValue(item_id-1, int(other_user[0])-1)}
		// Append in to normalized_rating_slice
		normalized_rating_slice = append(normalized_rating_slice, normalize_rating)
	}
	// Calculate and return prediction value
	return CalculatePredictionValue(similarity_slice, normalized_rating_slice)
}

// Find users who rate item i, sort them on decrease order of giving rate values
func (self *MatrixSlice) FindUsersWhoRateItem(item_id float64, user_id float64) [][]float64 {
	var users_slice [][]float64
	for i := 0; i < self.GetNumberOfRows(); i++ {
		if self.GetMatrixValue(i, 1) == item_id && self.GetMatrixValue(i, 0) != user_id {
			users_slice = append(users_slice, append([]float64{},
				self.GetMatrixValue(i, 0),
				self.GetMatrixValue(i, 2)))
		}
	}
	sort.Slice(users_slice, func(i, j int) bool {
		return users_slice[i][1] > users_slice[j][1]
	})
	return users_slice
}

// Select first n-th of slice
func SelectFirstElementsOfSlice(input_slice [][]float64, n int) [][]float64 {
	if n > len(input_slice) {
		log.Printf("The number of elements to take is larger than the length of input slice")
		n = len(input_slice)
	}
	var output_slice [][]float64
	for i := 0; i < n; i++ {
		output_slice = append(output_slice, input_slice[i])
	}
	return output_slice
}

// Get user vector
func (self *MatrixSlice) GetUserVector(user_id float64) []float64 {
	column_data := self.GetColumn(int(user_id - 1))
	return column_data
}

// Calculate the prediction value
func CalculatePredictionValue(similarity_slice [][]float64, normalized_rating_slice [][]float64) float64 {
	var numerator, denominator float64
	for i := 0; i < len(similarity_slice); i++ {
		numerator += similarity_slice[i][1] * normalized_rating_slice[i][1]
		denominator += math.Abs(similarity_slice[i][1])
	}
	return numerator / denominator
}

// Recommend item to user
func (self *CF) Recommend(user_id int) [][]float64 {
	var rated_item_ids []float64
	for i := 0; i < self.Y_data.GetNumberOfRows(); i++ {
		if self.Y_data.data[i][0] == float64(user_id) {
			rated_item_ids = append(rated_item_ids, self.Y_data.data[i][1])
		}
	}

	fmt.Printf("rated_item_ids: %v\n", rated_item_ids)

	var db_items []products.Product
	db := common.GetDB()
	db.Where("products.is_deleted = false").Find(&db_items)

	var unrated_item_ids []float64

	for _, db_item := range db_items {
		if !contains(rated_item_ids, float64(db_item.ID)) {
			unrated_item_ids = append(unrated_item_ids, float64(db_item.ID))
		}
	}

	fmt.Printf("unrated_item_ids: %v\n", unrated_item_ids)
	// Create slice to store item id and predict value, has form [[unrated_item_id predict_value] ...]
	var predict_item_slice [][]float64
	for _, unrated_item_id := range unrated_item_ids {
		predict_item := []float64{unrated_item_id,
			self.Predict(user_id, int(unrated_item_id))}
		predict_item_slice = append(predict_item_slice, predict_item)
	}
	fmt.Printf("predict_item_slice: %v\n", predict_item_slice)
	// Recommend item that has predict value > 0
	var recommend_items_slice [][]float64
	for _, i := range predict_item_slice {
		if i[1] > 0 {
			recommend_items_slice = append(recommend_items_slice, i)
		}
	}

	fmt.Printf("recommend_items_slice: %v\n", recommend_items_slice)
	return recommend_items_slice
}

// Process vector with zero value (which causes error when calculating)
// cosine similarity) by add 1e-8 (small value)
func ProcessIfZeroVector(vector []float64) {
	for _, i := range vector {
		if i != 0 {
			return
		}
	}

	vector[0] += 1e-8
}

// Check if slice contain element
func contains(slice []float64, element float64) bool {
	for _, v := range slice {
		if v == element {
			return true
		}
	}
	return false
}

// Refresh matrix when have new row
// TODO

// Content-based recommendation

// Get category vectors, has form [category_1 category_2]
func GetCategoryVectors() [][]string {
	var vector [][]string
	var categories []categories.Category

	db := common.GetDB()
	db.Find(&categories)

	for _, category := range categories {
		vector = append(vector, []string{strconv.Itoa(int(category.ID)), category.Name})
	}
	return vector
}

// Create item features matrix, each row is a feature vector and has form [[product_id 0 0 1 0 1] ...]
func GetItemFeaturesMatrix() MatrixSlice {
	var products_list []products.Product
	categories_vector := GetCategoryVectors()

	db := common.GetDB()
	db.Find(&products_list)

	// This matrix has form [[product_id 0 1 0], [product_id 0 1 0]]
	matrix := NewMatrix()

	// Number of row for defining a gonum library matrix

	// Loop through each product
	for _, product := range products_list {
		// Row to insert into matrix, has form [product_id 0 1 0]
		var matrix_row []float64

		// First element of row will be product id
		matrix_row = append(matrix_row, float64(product.ID))

		// Loop through each category
		for _, category_id := range categories_vector {
			fmt.Printf("category_id: %v\n", category_id)
			//  Row to insert data from gorm
			var category_product products.CategoryProductRelation
			result := db.Where("category_product_relations.category_id = ? AND category_product_relations.product_id = ?", category_id[0], product.ID).Find(&category_product)
			if result.RowsAffected == 0 {
				matrix_row = append(matrix_row, 0)
				continue
			}
			matrix_row = append(matrix_row, 1)
		}

		matrix.InsertRow(matrix_row)
	}

	return matrix
}

// Get rated items profile matrix
func GetRatedItemFeaturesMatrix(rated_matrix [][]float64) *mat.Dense {
	// Gonum matrix
	lib_matrix := mat.NewDense(len(rated_matrix), 2, nil)
	fmt.Printf("len(rated_matrix): %v\n", len(rated_matrix))

	// Get feature matrix
	feature_matrix := GetItemFeaturesMatrix()

	for i, rated_matrix_element := range rated_matrix {
		for _, feature_matrix_element := range feature_matrix.data {
			if rated_matrix_element[0] == feature_matrix_element[0] {
				fmt.Printf("i: %v\n", i)
				lib_matrix.SetRow(i, []float64{feature_matrix_element[1], feature_matrix_element[2]})
			}
		}
	}

	return lib_matrix
}

// Get products rated by user, return a matrix has form [[product_id rating] ....]
func GetProductsRatedByUser(user_id uint) ([][]float64, *mat.Dense) {
	var rated_matrix [][]float64

	db := common.GetDB()
	var ratings []products.Rating

	db.Find(&ratings)

	for _, rating := range ratings {
		if rating.UserID == uint64(user_id) {
			rated_matrix = append(rated_matrix, []float64{float64(rating.ProductID),
				float64(rating.Rate)})
		}
	}

	lib_matrix := mat.NewDense(len(rated_matrix), 1, nil)

	for i, rated := range rated_matrix {
		lib_matrix.SetRow(i, []float64{rated[1]})
	}

	return rated_matrix, lib_matrix
}
