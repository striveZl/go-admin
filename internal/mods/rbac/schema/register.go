package schema

import (
	"time"
)

type RegisterRequest struct {
	Nickname string `json:"nickname" form:"nickname" binding:"required,min=1,max=32"`
	Password string `json:"password" form:"password" binding:"required,min=8,max=72"`
	Email    string `json:"email" form:"email" binding:"required,email"`
	Code     string `json:"code" form:"code" binding:"required,len=6,numeric"`
}

type RegisterUser struct {
	ID           int64     `gorm:"column:id;primaryKey"`
	Nickname     string    `gorm:"column:nickname"`
	PasswordHash string    `gorm:"column:password_hash"`
	Email        string    `gorm:"column:email"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

type RegisterID struct {
	ID int64 `json:"id"`
}

type RegisterResponse struct {
	Success bool       `json:"success"`
	Data    RegisterID `json:"data"`
}
