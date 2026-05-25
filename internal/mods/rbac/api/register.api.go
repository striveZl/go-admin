package api

import (
	"fmt"
	"go-admin/internal/mods/rbac/biz"
	"go-admin/internal/mods/rbac/schema"
	"go-admin/pkg/errors"
	"go-admin/pkg/util"

	"github.com/gin-gonic/gin"
)

type Register struct {
	registerBIZ *biz.RegisterBIZ
}

func NewRegister(registerBIZ *biz.RegisterBIZ) *Register {
	return &Register{registerBIZ: registerBIZ}
}

func (register *Register) PostRegister(c *gin.Context) {
	if register == nil || register.registerBIZ == nil {
		util.ResError(c, errors.InternalServerError("", "register biz is not initialized"))
		return
	}

	var req schema.RegisterRequest

	if err := c.ShouldBind(&req); err != nil {
		util.ResError(c, errors.BadRequest("", "invalid register params"), 400)
		return
	}

	data, err := register.registerBIZ.Register(c.Request.Context(), &req)

	if err != nil {
		util.ResError(c, err)
		return
	}

	util.ResSuccess(c, data)

}

func (register *Register) PostCaptcha(c *gin.Context) {
	if register == nil || register.registerBIZ == nil {
		util.ResError(c, errors.InternalServerError("", "Captcha biz is not initialized"))
		return
	}

	var req schema.CaptchaRequest

	if err := c.ShouldBind(&req); err != nil {
		fmt.Println("邮箱:", err)
		fmt.Println("传入参数", req)
		util.ResError(c, errors.BadRequest("", "params error : %v", err), 400)
		return
	}

	data, err := register.registerBIZ.Captcha(c.Request.Context(), &req)

	if err != nil {
		util.ResError(c, err)
		return
	}

	util.ResSuccess(c, data)

}
