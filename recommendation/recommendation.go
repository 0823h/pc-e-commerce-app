package recommendation

import (
	"fmt"
	"log"
	"math"
	"sort"
	"tmdt-backend/common"
	"tmdt-backend/products"
	"tmdt-backend/users"

	"github.com/gaspiman/cosine_similarity"
)

type MatrixSlice struct {
	data [][]float64
}

// Creating a new matrix
func NewMatrix() *MatrixSlice {
	matrix := &MatrixSlice{nil}
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
	Y_data    *MatrixSlice
	Ybar_data *MatrixSlice
	n_users   int
	n_items   int
	k         int
}

// Creating a new CF
func NewCF() CF {
	y_data := NewMatrix()
	y_data.LoadRating()
	ybar_data := NewMatrix()

	for i := 0; i < y_data.GetNumberOfRows(); i++ {
		ybar_data.data = append(ybar_data.data, y_data.data[i])
	}

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
	mean_users := cf.Y_data.CalculateMeanFromMatrix()
	cf.Normalize_Y()
	fmt.Printf("cf.Y_data.data: %v\n", cf.Y_data.data)
	fmt.Printf("cf.Ybar_data.data: %v\n", cf.Ybar_data.data)
	fmt.Printf("mean_users: %v\n", mean_users)
	fmt.Printf("y_matrix: %v\n", cf.Ybar_data)
	cf.FormMatrix()
	fmt.Printf("cf.Ybar_data.data: %v\n", cf.Ybar_data.data)
	user_vector := cf.Ybar_data.GetUserVector(1)
	fmt.Printf("user_vector: %v\n", user_vector)
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
	distant, err := cosine_similarity.Cosine(user_1, user_2)
	if err != nil {
		log.Panic("Calculate distant fail!")
	}
	return distant
}

// Prediction function
func (self *CF) Predict(user_id int, item_id int) float64 {
	// Get ids of all users who rate the item, has form [[other_user_id rating_on_item], ...]
	users_slice := self.Y_data.FindUsersWhoRateItem(float64(item_id), float64(user_id))
	// Get vector of user for cosine similarity, has form [[rating_on_item1 rating_on_item2 rating_on_item3 ...]]
	user_vector := self.Ybar_data.GetUserVector(float64(user_id))
	// Slice to store all calculated similarity, has form [[other_user_id similarity], ...]
	var similarity_slice [][]float64
	// Calculate similarity and append into the slice above
	for _, other_user := range users_slice {
		// Get vector of other user, has form [[rating_on_item1 rating_on_item2 rating_on_item3 ...]]
		other_user_vector := self.Ybar_data.GetUserVector(other_user[0])
		// Creating a temporary slice for storing one calculated similarity above, has form [other_user_id similarity]
		temp_slice := []float64{other_user[0], CalculateSimilarity(user_vector, other_user_vector)}
		// Append into similarity_slice
		similarity_slice = append(similarity_slice, temp_slice)
	}

	// Sort similarity_slice following decrease order of similarity
	sort.Slice(similarity_slice, func(i, j int) bool {
		return similarity_slice[i][1] > similarity_slice[j][1]
	})

	// Select first n-th elements which is first n-th other-user that is most similar to user (n is k specified in struct CF)
	similarity_slice = SelectFirstElementsOfSlice(similarity_slice, self.k)

	// Get normalized rating of them on item
	// Create a slice to store normalize_rating, has form [[other_user_id normalized_rating], ...]
	var normalized_rating_slice [][]float64
	// Loop through similarity slice to get other user id and get normalize rating
	for _, other_user := range similarity_slice {
		// Get normalize rating of other user, has form [other_user_id normalized_rating]
		normalize_rating := []float64{other_user[1],
			self.Ybar_data.GetMatrixValue(item_id, int(other_user[1]))}
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
		log.Fatalf("The number of elements to take is larger than the length of input slice")
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
