package util

import "go-admin/pkg/errors"

const (
	reqBodyKey        = "req-body"
	resBodyKey        = "res-body"
	TreePathDelimiter = "."
)

type ResponseResult struct {
	Success bool          `json:"success"`
	Data    interface{}   `json:"data,omitempty"`
	Total   int64         `json:"total,omitempty"`
	Error   *errors.Error `json:"error,omitempty"`
}
