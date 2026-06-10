package dal

import (
	"context"
	"go-admin/internal/mods/rbac/schema"
	"go-admin/pkg/errors"

	"gorm.io/gorm"
)

type UserDAL struct {
	db *gorm.DB
}

func NewUser(db *gorm.DB) *UserDAL {
	return &UserDAL{db: db}
}

func (d *UserDAL) ExistsByAuthEmail(ctx context.Context, userEmail string) (bool, error) {
	if d == nil || d.db == nil {
		return false, errors.InternalServerError("", "database is not initialized")
	}

	var count int64

	//使用model可以自动支持软删除
	err := d.db.WithContext(ctx).Model(&schema.UserAuth{}).Where("identifier=? AND auth_type=?", userEmail, schema.AuthTypeEmail).Count(&count).Error

	if err != nil {
		return false, errors.InternalServerError("", "query email: %v", err)
	}

	return count > 0, nil

}
