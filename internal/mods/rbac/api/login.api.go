package api

import (
	"go-admin/internal/mods/rbac/biz"
	"go-admin/pkg/errors"
	"go-admin/pkg/util"

	"github.com/gin-gonic/gin"
)

type Login struct {
	loginBIZ *biz.Login
}

func NewLogin(loginBIZ *biz.Login) *Login {
	return &Login{loginBIZ: loginBIZ}
}

// 验证码
// @Summary 获取验证码 ID
// @Description 返回登录流程所需的验证码标识。
// @Tags RBAC
// @Produce json
// @Success 200 {object} schema.CaptchaResponse
// @Failure 500 {object} schema.ErrorResponse
// @Router /api/v1/captcha/id [get]
func (a *Login) GetCaptcha(c *gin.Context) {
	if a == nil || a.loginBIZ == nil {
		util.ResError(c, errors.InternalServerError("", "login biz is not initialized"))
		return
	}

	ctx := c.Request.Context()
	data, err := a.loginBIZ.GetCaptcha(ctx)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, data)

}
