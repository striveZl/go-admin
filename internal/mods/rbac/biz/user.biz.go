package biz

import (
	"context"
	"go-admin/internal/mods/rbac/dal"
	"go-admin/internal/mods/rbac/schema"
	pkgerrors "go-admin/pkg/errors"

	"gorm.io/gorm"
)

type User struct {
	dal *dal.UserDAL
}

func NewUser(userDAL *dal.UserDAL) *User {
	return &User{dal: userDAL}
}

func (a *User) GetByID(ctx context.Context, id uint) (*schema.User, error) {
	user, err := a.dal.GetByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, pkgerrors.NotFound("", "user not found")
		}
		return nil, err
	}
	return user, nil
}
