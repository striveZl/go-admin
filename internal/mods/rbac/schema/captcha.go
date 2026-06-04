package schema

type CaptchaRequest struct {
	Email string `json:"email" form:"email" binding:"required,email"`
}

type CaptchaResponse struct {
	ExpireSeconds int `json:"expire_seconds"`
}
