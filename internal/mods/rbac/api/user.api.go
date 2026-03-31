package api

import (
	"strconv"

	"go-admin/internal/mods/rbac/biz"
	"go-admin/pkg/errors"
	"go-admin/pkg/util"

	"github.com/gin-gonic/gin"
)

type User struct {
	userBIZ *biz.User
}

func NewUser(userBIZ *biz.User) *User {
	return &User{userBIZ: userBIZ}
}

// 获取用户详情
// @Summary 按 ID 获取用户
// @Description 根据主键从 users 表中返回一条用户记录。
// @Tags RBAC
// @Produce json
// @Param id path int true "用户 ID"
// @Success 200 {object} schema.UserResponse
// @Failure 400 {object} schema.ErrorResponse
// @Failure 404 {object} schema.ErrorResponse
// @Failure 500 {object} schema.ErrorResponse
// @Router /api/v1/users/{id} [get]
func (a *User) GetByID(c *gin.Context) {
	if a == nil || a.userBIZ == nil {
		util.ResError(c, errors.InternalServerError("", "user biz is not initialized"))
		return
	}

	rawID := c.Param("id")
	id, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil || id == 0 {
		util.ResError(c, errors.InternalServerError("", "invalid user id"), 400)
		return
	}

	user, err := a.userBIZ.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if e, ok := errors.As(err); ok && e != nil && e.Code == 404 {
			util.ResError(c, err, 404)
			return
		}
		util.ResError(c, err)
		return
	}

	util.ResSuccess(c, user)
}
