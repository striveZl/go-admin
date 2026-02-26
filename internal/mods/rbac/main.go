package rbac

import (
	"context"
	"fmt"
	"go-admin/internal/mods/rbac/api"

	"github.com/gin-gonic/gin"
)

type RBAC struct {
	Casbinx  *Casbinx
	LoginAPI *api.Login
}

func (a *RBAC) RegisterV1Routers(_ context.Context, v1 *gin.RouterGroup) error {
	if a == nil || a.LoginAPI == nil {
		return fmt.Errorf("login api is not initialized")
	}

	captcha := v1.Group("captcha")
	{
		captcha.GET("id", a.LoginAPI.GetCaptcha)
		// captcha.GET("image", a.LoginAPI.ResponseCaptcha)
	}

	return nil
}

func (a *RBAC) Release(ctx context.Context) error {
	if err := a.Casbinx.Release(ctx); err != nil {
		return err
	}
	return nil
}
