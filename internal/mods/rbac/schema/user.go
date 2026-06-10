package schema

import (
	"time"

	"gorm.io/gorm"
)

type GenderType string

const (
	GenderTypeMale   GenderType = "male"
	GenderTypeFemale GenderType = "female"
	GenderTypeOther  GenderType = "other"
)

type User struct {
	ID           int64          `gorm:"column:id;primaryKey"`
	Email        *string        `gorm:"column:email"`
	Nickname     string         `gorm:"column:nickname"`
	PasswordHash *string        `gorm:"column:password_hash"`
	Gender       *GenderType    `gorm:"column:gender;type:gender_type"`
	Birth        *time.Time     `gorm:"column:birth;type:date"`
	Avatar       *string        `gorm:"column:avatar"`
	Phone        *string        `gorm:"column:phone"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at"`
	CreatedAt    time.Time      `gorm:"column:created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at"`
}

func (User) TableName() string {
	return "users"
}
