package models

import (
	"time"
)

type OrderItem struct {
	ID        string    `json:"id" gorm:"not null;uniqueIndex;primary_key"`
	OrderID   string    `json:"order_id" gorm:"not null;index"`
	Order     Order     `gorm:"foreignKey:OrderID" json:"order"`
	ProductID string    `json:"product_id" gorm:"not null;index"`
	Product   Product   `gorm:"foreignKey:ProductID" json:"product"`
	Quantity  uint32    `json:"quantity" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type OrderItemResponse struct {
	ID        string    `json:"id"`
	OrderID   string    `json:"order_id"`
	ProductID string    `json:"product_id"`
	Product   Product   `json:"product"`
	Quantity  uint32    `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OrderItemCreateUpdate struct {
	OrderID   string `json:"order_id"`
	ProductID string `json:"product_id"`
	Quantity  uint32 `json:"quantity" validate:"required"`
}

type OrderItemDto struct {
	ProductID string `json:"product_id" validate:"required"`
	Quantity  uint32 `json:"quantity" validate:"required,min=1"`
}

func ToOrderItemResponse(orderItem OrderItem) OrderItemResponse {
	return OrderItemResponse{
		ID:        orderItem.ID,
		OrderID:   orderItem.OrderID,
		ProductID: orderItem.ProductID,
		Product:   orderItem.Product,
		Quantity:  orderItem.Quantity,
		CreatedAt: orderItem.CreatedAt,
		UpdatedAt: orderItem.UpdatedAt,
	}
}

func ToOrderItemResponses(orderItems []OrderItem) []OrderItemResponse {
	var responses []OrderItemResponse
	for _, item := range orderItems {
		responses = append(responses, ToOrderItemResponse(item))
	}
	return responses
}
