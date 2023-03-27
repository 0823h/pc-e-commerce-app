package recommendation

import (
	"fmt"
	"log"
	"strconv"
	"tmdt-backend/common"
	"tmdt-backend/products"

	"gonum.org/v1/gonum/mat"
)

type Matrix struct {
	matrix *mat.Dense
	row    int
	col    int
}

func NewMatrix(row int, col int, data []float64) Matrix {
	matrix := Matrix{mat.NewDense(row, col, data), row - 1, col}
	if row == 1 {
		row -= 1
	}
	return matrix
}

// func main() {
// 	matrix := NewMatrix(1, 3, nil)
// 	matrix.InsertRow(1, 2, 3)
// 	matrix.InsertRow(4, 5, 6)
// 	matrix.InsertRow(7, 8, 9)
// 	// matrix.InsertColumn(1, 2, 3)
// 	// matrix.InsertColumn(1, 2, 3)
// 	// matrix.InsertColumn(1, 2, 3, 4)
// 	// matrix.InsertColumn(1, 2, 3, 4)

// 	fmt.Printf("matrix: %v\n", matrix.matrix)
// 	// matrix.row = 1
// }

func (self *Matrix) InsertRow(elems ...float64) {
	if len(elems) != self.col {
		log.Fatalf("elems must have length: " + strconv.Itoa(self.col))
	}
	self.row += 1
	if self.row == 1 {
		self.matrix = mat.NewDense(self.row, self.col, elems)
		// fmt.Println(self.row)
		return
	}
	self.matrix = mat.NewDense(self.row, self.col, append(self.matrix.RawMatrix().Data, elems...))
}

// func (self *Matrix) InsertColumn(elems ...float64) {
// 	if len(elems) != self.row {
// 		log.Fatalf("elems must have length: " + strconv.Itoa(self.row))
// 	}
// 	self.col += 1

// 	var new_data_slice []float64

// 	fmt.Printf("len(self.matrix.RawMatrix().Data): %v\n", len(self.matrix.RawMatrix().Data))
// 	for i := 0; i < len(self.matrix.RawMatrix().Data); i++ {
// 		if i%self.row == 0 && i != 0 {
// 			new_data_slice = append(new_data_slice, 0)
// 			// continue
// 		}
// 		new_data_slice = append(new_data_slice, self.matrix.RawMatrix().Data[i])
// 	}
// 	new_data_slice = append(new_data_slice, 0)

// 	fmt.Println("hello")
// 	fmt.Println(new_data_slice)
// 	fmt.Printf("self.row: %v\n", self.row)
// 	fmt.Printf("self.col: %v\n", self.col)
// 	new_matrix := NewMatrix(self.row, self.col, new_data_slice)
// 	fmt.Printf("new_matrix: %v\n", new_matrix)

// 	// n ew_matrix.matrix.SetCol(self.col-1, data_slice)
// 	fmt.Printf("new_matrix: %v\n", new_matrix.matrix)
// 	self = &new_matrix
// }

func (self *Matrix) LoadRatings() {
	db := common.GetDB()
	var ratings []products.Rating
	db.Find(&ratings)
	for _, rating := range ratings {
		self.InsertRow(float64(rating.UserID), float64(rating.ProductID), float64(rating.Rate))
	}
}

func InitMatrix() {
	matrix := NewMatrix(1, 3, nil)
	matrix.LoadRatings()
	fmt.Printf("matrix: %v\n", matrix)
}
