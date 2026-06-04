package api

import (
	"go-admin/internal/mods/rbac/biz"
	"go-admin/internal/mods/rbac/schema"
	"go-admin/pkg/errors"
	"go-admin/pkg/util"

	"github.com/gin-gonic/gin"
)

type CaptchaAPI struct {
	captchaBIZ *biz.CaptchaBIZ
}

func NewCaptcha(captchaBIZ *biz.CaptchaBIZ) *CaptchaAPI {
	return &CaptchaAPI{captchaBIZ: captchaBIZ}

}

func (captchaAPI *CaptchaAPI) PostEmailCaptcha(c *gin.Context) {
	if captchaAPI == nil || captchaAPI.captchaBIZ == nil {
		util.ResError(c, errors.InternalServerError("", "Captcha biz is not initialized"))
		return
	}

	var req schema.CaptchaRequest

	if err := c.ShouldBind(&req); err != nil {
		util.ResError(c, errors.BadRequest("", "params error : %v", err), 400)
		return
	}

	data, err := captchaAPI.captchaBIZ.Captcha(c.Request.Context(), &req)

	if err != nil {
		util.ResError(c, err)
		return
	}

	util.ResSuccess(c, data)

}
