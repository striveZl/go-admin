package wirex

import (
	"context"
	"go-admin/internal/config"
	"go-admin/internal/mods"
	rbac2 "go-admin/internal/mods/rbac"
	"go-admin/internal/mods/rbac/api"
	"go-admin/internal/mods/rbac/biz"
	"go-admin/internal/mods/rbac/dal"
	"go-admin/pkg/email"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Mods struct {
	RBAC *rbac2.RBAC
}

func BuildInjector(_ context.Context, db *gorm.DB, redisClient *redis.Client) (*Injector, func(), error) {

	injector := &Injector{
		M: &mods.Mods{},
	}
	clearFn := func() {}

	casbinx := &rbac2.Casbinx{}

	emailSender := email.NewSMTPSender(
		config.C.Email.EmailApiKey,
		config.C.Email.From,
	)

	userDAL := dal.NewUser(db)

	captchaBIZ := biz.NewCaptcha(userDAL, redisClient, emailSender, config.C.Secret.CaptchaSecret)
	captchaAPI := api.NewCaptcha(captchaBIZ)

	registerDAL := dal.NewRegister(db)
	registerBIZ := biz.NewRegister(registerDAL, userDAL, captchaBIZ)
	registerAPI := api.NewRegister(registerBIZ)

	rbacRBAC := &rbac2.RBAC{
		Casbinx:     casbinx,
		RegisterAPI: registerAPI,
		CaptchaAPI:  captchaAPI,
	}

	modsMods := &mods.Mods{
		RBAC: rbacRBAC,
	}

	injector = &Injector{
		M: modsMods,
	}

	return injector, clearFn, nil
}
