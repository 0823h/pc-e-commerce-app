package products

import (
	"tmdt-backend/manufacturers"

	"github.com/gin-gonic/gin"
)

type ProductSerializer struct {
	C *gin.Context
	ProductModel
}

type ProductResponse struct {
	ID             uint                            `json:"-"`
	SKU            string                          `json:"sku"`
	Name           string                          `json:"name"`
	Description    string                          `json:"description"`
	Images         *string                         `json:"images"`
	Rating         float32                         `json:"rating"`
	Price          float64                         `json:"price0"`
	Quantity       uint                            `json:"quantity"`
	SoldAmount     uint                            `json:"sold_amount"`
	ManufacturerId uint                            `json:"manufacturer_id"`
	Manufacturer   manufacturers.ManufacturerModel `json:"manufacturer"`
}

func (self *ProductSerializer) Response() ProductResponse {
	product := ProductResponse{
		ID:             self.ID,
		SKU:            self.SKU,
		Name:           self.Name,
		Description:    self.Description,
		Images:         self.Images,
		Rating:         self.Rating,
		Price:          self.Price,
		Quantity:       self.Quantity,
		SoldAmount:     self.SoldAmount,
		ManufacturerId: self.ManufacturerId,
		Manufacturer:   self.Manufacturer,
	}
	return product
}
