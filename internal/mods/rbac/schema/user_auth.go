package schema

import (
	"time"

	"gorm.io/gorm"
)

type AuthType string

const (
	AuthTypeEmail  AuthType = "email"
	AuthTypeGoogle AuthType = "google"
	AuthTypePhone  AuthType = "phone"
)

type UserAuth struct {
	ID         int64          `gorm:"column:id;primaryKey"`
	UserID     int64          `gorm:"column:user_id"`
	AuthType   AuthType       `gorm:"column:auth_type;type:auth_type"`
	CreatedAt  time.Time      `gorm:"column:created_at"`
	Identifier string         `gorm:"column:identifier"`
	Credential *string        `gorm:"column:credential"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (UserAuth) TableName() string {
	return "user_auth"
}
