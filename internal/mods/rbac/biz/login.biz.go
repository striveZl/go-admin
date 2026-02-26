package biz

import (
	"context"
	"go-admin/internal/config"

	"github.com/LyricTian/captcha"

	"go-admin/internal/mods/rbac/schema"
)

type Login struct {
}

func (a *Login) GetCaptcha(_ context.Context) (*schema.Captcha, error) {
	return &schema.Captcha{
		CaptchaID: captcha.NewLen(config.C.Util.Captcha.Length),
	}, nil
}
