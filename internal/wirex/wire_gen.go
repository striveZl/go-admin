package wirex

import (
	"context"
	"go-admin/internal/mods"
	rbac2 "go-admin/internal/mods/rbac"
	"go-admin/internal/mods/rbac/api"
	"go-admin/internal/mods/rbac/biz"
)

type Mods struct {
	RBAC *rbac2.RBAC
}

func BuildInjector(_ context.Context) (*Injector, func(), error) {
	injector := &Injector{
		M: &mods.Mods{},
	}
	clearFn := func() {}

	casbinx := &rbac2.Casbinx{}
	loginBiz := &biz.Login{}
	loginAPI := api.NewLogin(loginBiz)

	rbacRBAC := &rbac2.RBAC{

		Casbinx:  casbinx,
		LoginAPI: loginAPI,
	}

	modsMods := &mods.Mods{
		RBAC: rbacRBAC,
	}

	injector = &Injector{
		M: modsMods,
	}

	return injector, clearFn, nil
}
