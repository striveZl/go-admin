package rbac

import (
	"context"
	"fmt"
	"go-admin/internal/mods/rbac/api"

	"github.com/gin-gonic/gin"
)

type RBAC struct {
	Casbinx     *Casbinx
	RegisterAPI *api.RegisterAPI
	CaptchaAPI  *api.CaptchaAPI
}

func (a *RBAC) RegisterV1Routers(_ context.Context, v1 *gin.RouterGroup) error {
	if a == nil {
		return fmt.Errorf("rbac apis are not initialized")
	}

	register := v1.Group("register")
	{
		register.POST("create", a.RegisterAPI.PostRegister)
	}

	captcha := v1.Group("captcha")
	{
		captcha.POST("email", a.CaptchaAPI.PostEmailCaptcha)
	}

	return nil
}

func (a *RBAC) Release(ctx context.Context) error {
	if err := a.Casbinx.Release(ctx); err != nil {
		return err
	}
	return nil
}
