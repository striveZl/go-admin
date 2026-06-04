package dal

import (
	"context"
	"go-admin/pkg/errors"

	"gorm.io/gorm"
)

type UserDAL struct {
	db *gorm.DB
}

func NewUser(db *gorm.DB) *UserDAL {
	return &UserDAL{db: db}
}

func (d *UserDAL) ExistsByEmail(ctx context.Context, userEmail string) (bool, error) {
	if d == nil || d.db == nil {
		return false, errors.InternalServerError("", "database is not initialized")
	}

	var count int64

	err := d.db.WithContext(ctx).Table("users").Where("email=?", userEmail).Count(&count).Error

	if err != nil {
		return false, errors.InternalServerError("", "query email: %v", err)
	}

	return count > 0, nil

}
