package products

import (
	"tmdt-backend/manufacturers"

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
	ID             string                     `json:"id"`
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
