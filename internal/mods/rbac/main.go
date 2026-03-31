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
	UserAPI  *api.User
}

func (a *RBAC) RegisterV1Routers(_ context.Context, v1 *gin.RouterGroup) error {
	if a == nil || a.LoginAPI == nil || a.UserAPI == nil {
		return fmt.Errorf("rbac apis are not initialized")
	}

	captcha := v1.Group("captcha")
	{
		captcha.GET("id", a.LoginAPI.GetCaptcha)
		// captcha.GET("image", a.LoginAPI.ResponseCaptcha)
	}

	users := v1.Group("users")
	
	{
		users.GET("/:id", a.UserAPI.GetByID)
	}

	return nil
}

func (a *RBAC) Release(ctx context.Context) error {
	if err := a.Casbinx.Release(ctx); err != nil {
		return err
	}
	return nil
}
