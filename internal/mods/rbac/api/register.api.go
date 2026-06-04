package api

import (
	"go-admin/internal/mods/rbac/biz"
	"go-admin/internal/mods/rbac/schema"
	"go-admin/pkg/errors"
	"go-admin/pkg/util"

	"github.com/gin-gonic/gin"
)

type RegisterAPI struct {
	registerBIZ *biz.RegisterBIZ
}

func NewRegister(registerBIZ *biz.RegisterBIZ) *RegisterAPI {
	return &RegisterAPI{registerBIZ: registerBIZ}
}

func (registerAPI *RegisterAPI) PostRegister(c *gin.Context) {

	if registerAPI == nil || registerAPI.registerBIZ == nil {
		util.ResError(c, errors.InternalServerError("", "register biz is not initialized"))
		return
	}

	var req schema.RegisterRequest
	if err := c.ShouldBind(&req); err != nil {
		util.ResError(c, errors.BadRequest("", "invalid register params"), 400)
		return
	}

	data, err := registerAPI.registerBIZ.Register(c.Request.Context(), &req)
	if err != nil {
		util.ResError(c, err)
		return
	}

	util.ResSuccess(c, data)

}
