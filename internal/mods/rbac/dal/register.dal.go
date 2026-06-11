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

func (d *RegisterDAL) RegisterCreate(ctx context.Context, user *schema.User, passwordHash *string) (*schema.RegisterID, error) {
	if d == nil || d.db == nil {
		return nil, errors.InternalServerError("", "register dal is not initialized")
	}

	if user == nil {
		return nil, errors.BadRequest("", "user is required")
	}

	if user.Email == nil || *user.Email == "" {
		return nil, errors.BadRequest("", "email is required")
	}

	//创建事务以确保两条sql能够都执行成功
	err := d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		//创建用户
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		user_auth := &schema.UserAuth{
			UserID:     user.ID,
			AuthType:   schema.AuthTypeEmail,
			Identifier: *user.Email,
			Credential: passwordHash,
		}

		//创建user_auth
		if err := tx.Create(user_auth).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, errors.InternalServerError("", "create user: %v", err)
	}

	return &schema.RegisterID{
		ID: user.ID,
	}, nil

}
