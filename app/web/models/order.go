package models

import (
	"time"
)

type Order struct {
	ID           string      `json:"id" gorm:"not null;uniqueIndex;primary_key"`
	UserID       string      `json:"user_id" gorm:"not null"`
	OrderItems   []OrderItem `json:"order_items" gorm:"foreignKey:OrderID"`
	IsPaid       bool        `json:"is_paid"`
	Status       string      `json:"status"`
	CustomerName string      `json:"customer_name"`
	Phone        string      `json:"phone"`
	TotalPrice   float64     `json:"total_price"`
	Address      string      `json:"address"`
	CreatedAt    time.Time   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time   `json:"updated_at" gorm:"autoUpdateTime"`
}

type OrderResponse struct {
	ID           string              `json:"id"`
	UserID       string              `json:"user_id"`
	OrderItems   []OrderItemResponse `json:"order_items"`
	IsPaid       bool                `json:"is_paid"`
	CustomerName string              `json:"customer_name"`
	Phone        string              `json:"phone"`
	Address      string              `json:"address"`
	TotalPrice   float64             `json:"total_price"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

type OrderCreateUpdate struct {
	IsPaid     bool      `json:"is_paid"`
	Phone      string    `json:"phone"`
	Address    string    `json:"address"`
	TotalPrice float64   `json:"total_price"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func ToOrderResponse(order Order) OrderResponse {
	var orderItems []OrderItemResponse

	for _, orderItem := range order.OrderItems {
		orderItems = append(orderItems, ToOrderItemResponse(orderItem))
	}
	return OrderResponse{
		ID:           order.ID,
		UserID:       order.UserID,
		OrderItems:   orderItems,
		IsPaid:       order.IsPaid,
		CustomerName: order.CustomerName,
		Phone:        order.Phone,
		Address:      order.Address,
		TotalPrice:   order.TotalPrice,
		CreatedAt:    order.CreatedAt,
		UpdatedAt:    order.UpdatedAt,
	}
}

func ToOrderResponses(orders []Order) []OrderResponse {
	var responses []OrderResponse

	for _, order := range orders {
		responses = append(responses, ToOrderResponse(order))
	}
	return responses
}
