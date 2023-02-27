package products

import (
	"tmdt-backend/manufacturers"
	"tmdt-backend/users"

	"github.com/gin-gonic/gin"
)

type ProductsSerializer struct {
	C        *gin.Context
	Products []Product
}

type ProductSerializer struct {
	C *gin.Context
	Product
}

type ProductResponse struct {
	ID             uint64                     `json:"id"`
	SKU            string                     `json:"sku"`
	Name           string                     `json:"name"`
	Description    string                     `json:"description"`
	Images         *string                    `json:"images"`
	Rating         float32                    `json:"rating"`
	Price          float64                    `json:"price"`
	Quantity       uint                       `json:"quantity"`
	SoldAmount     uint                       `json:"sold_amount"`
	ManufacturerID uint                       `json:"manufacturer_id"`
	Manufacturer   manufacturers.Manufacturer `json:"manufacturer"`
}

func (self *ProductSerializer) Response() ProductResponse {
	response := ProductResponse{
		ID:             self.ID,
		SKU:            self.SKU,
		Name:           self.Name,
		Description:    self.Description,
		Images:         self.Images,
		Rating:         self.Rating,
		Price:          self.Price,
		Quantity:       self.Quantity,
		SoldAmount:     self.SoldAmount,
		ManufacturerID: self.ManufacturerID,
		Manufacturer:   self.Manufacturer,
	}
	return response
}

func (self *ProductsSerializer) Response() []ProductResponse {
	response := []ProductResponse{}
	for _, product := range self.Products {
		serializer := ProductSerializer{self.C, product}
		response = append(response, serializer.Response())
	}
	return response
}

type RatingResponse struct {
	ID            uint64     `json:"id"`
	UserID        uint64     `json:"user_id"`
	User          users.User `json:"user"`
	ProductID     uint64     `json:"product_id"`
	Product       Product    `json:"product"`
	Rate          uint       `json:"rate"`
	NumberOfClick uint       `json:"number_of_click"`
}

type RatingsSerializer struct {
	C       *gin.Context
	Ratings []Rating
}

type RatingSerializer struct {
	C *gin.Context
	Rating
}

func (self *RatingSerializer) Response() RatingResponse {
	response := RatingResponse{
		ID:            self.ID,
		UserID:        self.UserID,
		User:          self.User,
		ProductID:     self.ProductID,
		Product:       self.Product,
		Rate:          self.Rate,
		NumberOfClick: self.NumberOfClick,
	}
	return response
}

func (self *RatingsSerializer) Response() []RatingResponse {
	response := []RatingResponse{}
	for _, rating := range self.Ratings {
		serializer := RatingSerializer{self.C, rating}
		response = append(response, serializer.Response())
	}
	return response
}
