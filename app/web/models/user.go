package models

import (
	"time"
)

type User struct {
	ID        string    `json:"id" gorm:"not null;uniqueIndex;primary_key"`
	Name      string    `json:"name" gorm:"not null;type:varchar(50)"`
	Email     string    `json:"email" gorm:"not null;unique;type:varchar(100)"`
	Password  string    `json:"password" gorm:"not null;type:varchar(100)"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserLoginResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type UserCreate struct {
	Name     string `validate:"required,min=4,max=50"`
	Email    string `validate:"required,email"`
	Phone    string `validate:"required,min=8,max=13"`
	Password string `validate:"required,min=6,max=50"`
	Address  string `validate:"required,min=6,max=100"`
}

type UserUpdate struct {
	Phone    string `validate:"required,min=8,max=13"`
	Password string `validate:"required,min=6,max=50"`
	Address  string `validate:"required,min=6,max=100"`
}

type UserLogin struct {
	Email    string `validate:"required,email" json:"email"`
	Password string `validate:"required" json:"password"`
}

func ToUserReponse(user User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		Address:   user.Address,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
