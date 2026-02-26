package util

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-admin/pkg/errors"
	"go-admin/pkg/logging"
	"go.uber.org/zap"
	"net/http"
)

func ResOK(c *gin.Context) {
	ResJSON(c, http.StatusOK, ResponseResult{
		Success: true,
	})
}

func ResJSON(c *gin.Context, status int, v interface{}) {
	buf, err := json.Marshal(v) //把 Go 里的数据结构转换成 JSON 字节流

	if err != nil {
		ctx := c.Request.Context()
		ctx = logging.NewTag(ctx, logging.TagKeySystem)
		logging.Context(ctx).Error("Failed to marshal response", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, ResponseResult{
			Error: errors.FromError(errors.InternalServerError("", "failed to marshal response")),
		})
		return
	}
	c.Set(resBodyKey, buf)
	c.Data(status, "application/json; charset=utf-8", buf)
	c.Abort() //终止后续中间件

}

func ResError(c *gin.Context, err error, status ...int) {
	if err == nil {
		err = errors.InternalServerError("", "unknown error")
	}

	var ierr *errors.Error
	if e, ok := errors.As(err); ok {
		ierr = e
	} else {
		ierr = errors.FromError(errors.InternalServerError("", "%s", err.Error()))
	}
	code := http.StatusInternalServerError
	if len(status) > 0 {
		code = status[0]
	}
	if code >= 500 {
		ctx := c.Request.Context()
		ctx = logging.NewTag(ctx, logging.TagKeySystem)
		ctx = logging.NewStack(ctx, fmt.Sprintf("%+v", err))
		logging.Context(ctx).Error("Internal server error", zap.Error(err))
		ierr.Detail = http.StatusText(http.StatusInternalServerError)
	}

	ierr.Code = int32(code)
	ResJSON(c, code, ResponseResult{Error: ierr})
}

func ResSuccess(c *gin.Context, v interface{}) {
	ResJSON(c, http.StatusOK, ResponseResult{
		Success: true,
		Data:    v,
	})
}
