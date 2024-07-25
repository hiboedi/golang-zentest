package models

import (
	"time"
)

type Image struct {
	ID        string    ` json:"id" gorm:"not null;uniqueIndex;primary_key"`
	ProductID string    `json:"product_id" gorm:"not null;index"`
	Product   Product   `gorm:"foreignKey:ProductID" json:"product"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type ImageResponse struct {
	ID        string    `json:"id"`
	ProductID string    `json:"product_id"`
	Product   Product   `json:"product"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ImageResponseHiddenProduct struct {
	ID        string    `json:"id"`
	ProductID string    `json:"product_id"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ImageCreate struct {
	URL string `json:"url" validate:"required"`
}

type ImageUpdate struct {
	ProductID string `json:"product_id"`
	URL       string `json:"url" validate:"required"`
}

func ToImageResponse(image Image) ImageResponse {
	return ImageResponse{
		ID:        image.ID,
		ProductID: image.ProductID,
		URL:       image.URL,
		CreatedAt: image.CreatedAt,
		UpdatedAt: image.UpdatedAt,
	}
}

func ToImageResponses(images []Image) []ImageResponse {
	var responses []ImageResponse

	for _, image := range images {
		responses = append(responses, ToImageResponse(image))
	}
	return responses
}
