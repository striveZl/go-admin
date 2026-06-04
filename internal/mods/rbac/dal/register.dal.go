package dal

import (
	"context"
	"go-admin/internal/mods/rbac/schema"
	"go-admin/pkg/errors"

	"gorm.io/gorm"
)

type RegisterDAL struct {
	db *gorm.DB
}

func NewRegister(db *gorm.DB) *RegisterDAL {
	return &RegisterDAL{db: db}
}

func (d *RegisterDAL) RegisterCreate(ctx context.Context, user *schema.RegisterUser) (*schema.RegisterID, error) {
	if d == nil || d.db == nil {
		return nil, errors.InternalServerError("", "register dal is not initialized")
	}

	err := d.db.WithContext(ctx).Table("users").Create(user).Error

	if err != nil {
		return nil, errors.InternalServerError("", "create user: %v", err)
	}

	return &schema.RegisterID{
		ID: user.ID,
	}, nil

}
