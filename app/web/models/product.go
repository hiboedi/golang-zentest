package models

import (
	"time"
)

type Product struct {
	ID        string    `json:"id" gorm:"not null;uniqueIndex;primary_key"`
	Category  string    `json:"category"`
	Name      string    `json:"name"`
	Price     float64   `json:"price" `
	Stock     uint32    `json:"stock"`
	Images    []Image   `gorm:"foreignKey:ProductID" json:"images"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type ProductResponse struct {
	ID        string    `json:"id"`
	Category  string    `json:"category"`
	Name      string    `json:"name"`
	Stock     uint32    `json:"stock"`
	Images    []Image   `json:"images"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProductCreateUpdate struct {
	Category string  `json:"category" validate:"required,min=4,max=50"`
	Name     string  `json:"name" validate:"required,min=4,max=50"`
	Price    float64 `json:"price" validate:"required,min=1"`
	Stock    uint32  `json:"stock" validate:"required,min=1"`
	Images   []Image `json:"images" alidate:"required"`
}

type ProductDto struct {
	Category string        `json:"category"`
	Name     string        `json:"name"`
	Price    float64       `json:"price"`
	Stock    uint32        `json:"stock"`
	Images   []ImageCreate `json:"images"`
}

func ToProductResponse(product Product) ProductResponse {
	return ProductResponse{
		ID:        product.ID,
		Category:  product.Category,
		Name:      product.Name,
		Stock:     product.Stock,
		Images:    product.Images,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}
}

func ToProductResponses(products []Product) []ProductResponse {
	var prodctResponse []ProductResponse

	for _, product := range products {
		prodctResponse = append(prodctResponse, ToProductResponse(product))
	}

	return prodctResponse
}
