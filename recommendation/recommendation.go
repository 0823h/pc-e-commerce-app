package recommendation

import (
	"fmt"
	"log"
	"tmdt-backend/common"
	"tmdt-backend/products"
)

// type Matrix struct {
// 	matrix *mat.Dense
// 	row    int
// 	col    int
// }

// func NewMatrix(row int, col int, data []float64) Matrix {
// 	matrix := Matrix{mat.NewDense(row, col, data), row - 1, col}
// 	if row == 1 {
// 		row -= 1
// 	}
// 	return matrix
// }

// // func main() {
// // 	matrix := NewMatrix(1, 3, nil)
// // 	matrix.InsertRow(1, 2, 3)
// // 	matrix.InsertRow(4, 5, 6)
// // 	matrix.InsertRow(7, 8, 9)
// // 	// matrix.InsertColumn(1, 2, 3)
// // 	// matrix.InsertColumn(1, 2, 3)
// // 	// matrix.InsertColumn(1, 2, 3, 4)
// // 	// matrix.InsertColumn(1, 2, 3, 4)

// // 	fmt.Printf("matrix: %v\n", matrix.matrix)
// // 	// matrix.row = 1
// // }

// func (self *Matrix) InsertRow(elems ...float64) {
// 	if len(elems) != self.col {
// 		log.Fatalf("elems must have length: " + strconv.Itoa(self.col))
// 	}
// 	self.row += 1
// 	if self.row == 1 {
// 		self.matrix = mat.NewDense(self.row, self.col, elems)
// 		// fmt.Println(self.row)
// 		return
// 	}
// 	self.matrix = mat.NewDense(self.row, self.col, append(self.matrix.RawMatrix().Data, elems...))
// }

// // func (self *Matrix) InsertColumn(elems ...float64) {
// // 	if len(elems) != self.row {
// // 		log.Fatalf("elems must have length: " + strconv.Itoa(self.row))
// // 	}
// // 	self.col += 1

// // 	var new_data_slice []float64

// // 	fmt.Printf("len(self.matrix.RawMatrix().Data): %v\n", len(self.matrix.RawMatrix().Data))
// // 	for i := 0; i < len(self.matrix.RawMatrix().Data); i++ {
// // 		if i%self.row == 0 && i != 0 {
// // 			new_data_slice = append(new_data_slice, 0)
// // 			// continue
// // 		}
// // 		new_data_slice = append(new_data_slice, self.matrix.RawMatrix().Data[i])
// // 	}
// // 	new_data_slice = append(new_data_slice, 0)

// // 	fmt.Println("hello")
// // 	fmt.Println(new_data_slice)
// // 	fmt.Printf("self.row: %v\n", self.row)
// // 	fmt.Printf("self.col: %v\n", self.col)
// // 	new_matrix := NewMatrix(self.row, self.col, new_data_slice)
// // 	fmt.Printf("new_matrix: %v\n", new_matrix)

// // 	// n ew_matrix.matrix.SetCol(self.col-1, data_slice)
// // 	fmt.Printf("new_matrix: %v\n", new_matrix.matrix)
// // 	self = &new_matrix
// // }

// func (self *Matrix) LoadRatings() {
// 	db := common.GetDB()
// 	var ratings []products.Rating
// 	db.Find(&ratings)
// 	for _, rating := range ratings {
// 		self.InsertRow(float64(rating.UserID), float64(rating.ProductID), float64(rating.Rate))
// 	}
// }

// func InitMatrix() {
// 	matrix := NewMatrix(1, 3, nil)
// 	matrix.LoadRatings()
// 	fmt.Printf("matrix.matrix: %v\n", matrix.matrix)
// }

// type CF struct {
// 	uuCF      bool
// 	Y_data    *Matrix
// 	k         int
// 	Ybar_data Matrix
// 	n_users   int64
// 	n_items   int64
// }

// func (self *CF) dist_func() {
// 	// TO DO
// }

// func NewCF(uuCF bool, Y_data *Matrix, k int) *CF {
// 	// Calculate number of users, number of items
// 	db := common.GetDB()
// 	var users []users.User
// 	var products []products.Product
// 	result := db.Where("products.is_deleted = ?", "false").Find(&products)
// 	n_items := result.RowsAffected

// 	result = db.Where("users.is_deleted = ?", "false").Find(&users)
// 	n_users := result.RowsAffected

// 	// Constructor
// 	var cf *CF
// 	cf.uuCF = uuCF
// 	cf.Y_data = Y_data
// 	// cf.Ybar_data =
// 	cf.n_items = n_items
// 	cf.n_users = n_users
// 	cf.k = k

// 	return cf
// }

// func (self *CF) normalize_Y() {
// 	// users := self.Y_data.matrix.ColView(0)
// 	self.Ybar_data = NewMatrix(self.Y_data.row, self.Y_data.col, nil)
// 	self.Ybar_data.matrix.Copy(self.Y_data.matrix)
// 	user_mean_vector := mat.NewVecDense(int(self.n_users), nil)
// 	var rating_slices []float64
// 	for i := int64(0); i < self.n_users; i++ {
// 		if self.Y_data.matrix.RawRowView()
// 	}

// }

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
	copy(ybar_data.data, y_data.data)

	n_users, _ := y_data.GetUsersLength()
	n_items := y_data.GetItemsLength()
	cf := CF{y_data, ybar_data, n_users, n_items, 2}
	return cf
}

// func (self *CF) Normalize_Y() {
// 	mean_users := self.Y_data.CalculateMeanFromMatrix()
// 	for i := 0; i <= self.Ybar_data.GetNumberOfRows(); i++ {
// 		for j := 0; j <= self.mean_users[i]
// 		self.Ybar_data[i] -=
// 	}
// }

// A main function for testing in Recommendation package
func InitMatrix() {
	cf := NewCF()
	mean_users := cf.Y_data.CalculateMeanFromMatrix()
	fmt.Printf("mean_users: %v\n", mean_users)
}
