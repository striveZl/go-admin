package schema

type Captcha struct {
	CaptchaID string `json:"captcha_id"`
}

type Login struct {
	Username string `json:"username"`
}
