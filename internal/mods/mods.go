package mods

import (
	"context"
	"go-admin/internal/mods/rbac"

	"github.com/gin-gonic/gin"
)

const apiPrefix = "/api/"

type Mods struct {
	RBAC *rbac.RBAC
}

func (a *Mods) RegisterRouters(ctx context.Context, e *gin.Engine) error {
	gAPI := e.Group(apiPrefix)
	v1 := gAPI.Group("/v1")

	if err := a.RBAC.RegisterV1Routers(ctx, v1); err != nil {
		return err
	}
	return nil
}

func (a *Mods) Release(ctx context.Context) error {
	if err := a.RBAC.Release(ctx); err != nil {
		return err
	}

	return nil
}
