package dal

import (
	"context"
	"errors"
	"fmt"
	"go-admin/internal/mods/rbac/schema"

	"gorm.io/gorm"
)

// UserDAL 封装用户查询相关的持久化依赖。
type UserDAL struct {
	db *gorm.DB
}

// NewUserDAL 使用注入的数据库句柄创建用户数据访问层。
func NewUserDAL(db *gorm.DB) *UserDAL {
	return &UserDAL{db: db}
}

// GetByID 按主键从 users 表查询单个用户。
func (r *UserDAL) GetByID(ctx context.Context, id uint) (*schema.User, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("user dal is not initialized")
	}

	var user schema.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("query user by id %d: %w", id, err)
	}

	return &user, nil
}
