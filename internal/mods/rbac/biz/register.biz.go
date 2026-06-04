package biz

import (
	"context"
	"go-admin/internal/mods/rbac/dal"
	"go-admin/internal/mods/rbac/schema"
	"go-admin/pkg/errors"
	"go-admin/pkg/util"
	"time"
)

type RegisterBIZ struct {
	registerDAL *dal.RegisterDAL
	userDAL     *dal.UserDAL
	captchaBIZ  *CaptchaBIZ
}

func NewRegister(registerDAL *dal.RegisterDAL, userDAL *dal.UserDAL,
	captchaBIZ *CaptchaBIZ) *RegisterBIZ {
	return &RegisterBIZ{registerDAL: registerDAL, userDAL: userDAL, captchaBIZ: captchaBIZ}
}

// 频率
const captchaCoolDown = time.Second * 60

// 验证码过期时间
const captchaTTL = time.Minute * 5

func (registerBIZ *RegisterBIZ) Register(ctx context.Context, req *schema.RegisterRequest) (*schema.RegisterID, error) {
	if registerBIZ == nil || registerBIZ.registerDAL == nil {
		return nil, errors.InternalServerError("", "register dal is not initialized")
	}

	existsByUserAccount, err := registerBIZ.userDAL.ExistsByEmail(ctx, req.Email)

	if err != nil {
		return nil, errors.InternalServerError("", "Email:%v", err)
	}

	if existsByUserAccount {
		return nil, errors.BadRequest("", "The Email already exists")
	}

	//验证验证码是否正确或过期

	if err := registerBIZ.captchaBIZ.verifyHashCaptchaCode(ctx, req.Email, req.Code); err != nil {
		return nil, err
	}

	passwordHash, err := util.HashPassword(req.Password)

	if err != nil {
		return nil, errors.InternalServerError("", "hash password:%v", err)
	}

	now := time.Now()

	user := &schema.RegisterUser{
		Nickname:     req.Nickname,
		PasswordHash: passwordHash,
		Email:        req.Email,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	registerId, err := registerBIZ.registerDAL.RegisterCreate(ctx, user)

	if err != nil {
		return nil, err
	}
	return registerId, err

}
