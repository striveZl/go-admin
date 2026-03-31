package db

import (
	"context"
	"database/sql"
	"fmt"
	"go-admin/internal/config"
	pkgerrors "go-admin/pkg/errors"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Store 封装 gorm 和 sql.DB 句柄，统一管理数据库生命周期。
type Store struct {
	gormDB *gorm.DB
	sqlDB  *sql.DB
}

// Open 根据运行时配置初始化 PostgreSQL 连接。
func Open(ctx context.Context, cfg config.Database) (*Store, error) {
	if cfg.Driver != "postgres" {
		return nil, pkgerrors.InternalServerError("", "unsupported database driver: %s", cfg.Driver)
	}
	if cfg.Host == "" {
		return nil, pkgerrors.InternalServerError("", "database host is required")
	}
	if cfg.User == "" {
		return nil, pkgerrors.InternalServerError("", "database user is required")
	}
	if cfg.DBName == "" {
		return nil, pkgerrors.InternalServerError("", "database name is required")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
		cfg.Timezone,
	)

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open postgres connection: %w", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("resolve sql db from gorm: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	if cfg.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return &Store{
		gormDB: gormDB,
		sqlDB:  sqlDB,
	}, nil
}

// Gorm 返回已初始化的 gorm 数据库句柄。
func (s *Store) Gorm() *gorm.DB {
	if s == nil {
		return nil
	}
	return s.gormDB
}

// SQL 返回已初始化的 sql.DB 句柄。
func (s *Store) SQL() *sql.DB {
	if s == nil {
		return nil
	}
	return s.sqlDB
}

// Close 关闭底层 sql.DB 连接。
func (s *Store) Close() error {
	if s == nil || s.sqlDB == nil {
		return nil
	}
	return s.sqlDB.Close()
}
