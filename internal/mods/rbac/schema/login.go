package schema

import pkgErrors "go-admin/pkg/errors"

type Captcha struct {
	CaptchaID string `json:"captcha_id"`
}

type Login struct {
	Username string `json:"username"`
}

type CaptchaResponse struct {
	Success bool    `json:"success"`
	Data    Captcha `json:"data"`
}

type ErrorResponse struct {
	Success bool             `json:"success"`
	Error   *pkgErrors.Error `json:"error,omitempty"`
}
