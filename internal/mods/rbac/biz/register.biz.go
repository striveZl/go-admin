package biz

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"go-admin/internal/mods/rbac/dal"
	"go-admin/internal/mods/rbac/schema"
	"go-admin/pkg/errors"
	"go-admin/pkg/util"
	"math/big"
	"time"

	"github.com/go-redis/redis/v8"
)

type RegisterBIZ struct {
	registerDAL   *dal.RegisterDAL
	redisClient   *redis.Client
	emailSender   EmailSender
	captchaSecret string
}

func NewRegister(registerBIZ *dal.RegisterDAL, redisClient *redis.Client, captchaSecret string, emailSender EmailSender) *RegisterBIZ {
	return &RegisterBIZ{registerDAL: registerBIZ, redisClient: redisClient, captchaSecret: captchaSecret, emailSender: emailSender}
}

// 频率
const captchaCoolDown = time.Second * 60

// 验证码过期时间
const captchaTTL = time.Minute * 5

func (registerBIZ *RegisterBIZ) Register(ctx context.Context, req *schema.RegisterRequest) (*schema.RegisterID, error) {
	if registerBIZ == nil || registerBIZ.registerDAL == nil {
		return nil, errors.InternalServerError("", "register dal is not initialized")
	}

	existsByUserAccount, err := registerBIZ.registerDAL.ExistsByEmail(ctx, req.Email)

	if err != nil {
		return nil, errors.InternalServerError("", "Email:%v", err)
	}

	if existsByUserAccount {
		return nil, errors.BadRequest("", "The Email already exists")
	}

	//验证验证码是否正确或过期

	if err := registerBIZ.verifyHashCaptchaCode(ctx, req.Email, req.Code); err != nil {
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

// 验证码发送
func (registerBIZ *RegisterBIZ) Captcha(ctx context.Context, req *schema.CaptchaRequest) (*schema.Captcha, error) {
	if registerBIZ == nil || registerBIZ.registerDAL == nil {
		return nil, errors.InternalServerError("", "Captcha dal is not initialized")
	}

	existsByEmail, err := registerBIZ.registerDAL.ExistsByEmail(ctx, req.Email)

	if err != nil {
		return nil, errors.InternalServerError("", "Email:%v", err)
	}

	if existsByEmail {
		return nil, errors.BadRequest("", "The Email already register")
	}

	if err := registerBIZ.checkCaptchaCoolDown(ctx, req.Email); err != nil {
		return nil, err
	}

	code, err := generateRandCaptcha()

	if err != nil {
		return nil, errors.InternalServerError("", "generate captcha: %v", err)
	}

	if err := registerBIZ.saveCaptcha(ctx, req.Email, code); err != nil {
		return nil, err
	}

	if err := registerBIZ.emailSender.SendCaptcha(ctx, req.Email, code); err != nil {
		registerBIZ.cleanupCaptchaOnSendFailed(ctx, req.Email)
		return nil, errors.InternalServerError("", "send captcha email: %v", err)
	}

	return &schema.Captcha{
		ExpireSeconds: int(captchaTTL.Seconds()),
	}, nil

}

// 限制发送频率
func (b *RegisterBIZ) checkCaptchaCoolDown(ctx context.Context, email string) error {
	coolDownKey := fmt.Sprintf("register:captcha:coolDown:%s", email)

	// 60秒内，无法再次发送
	ok, err := b.redisClient.SetNX(ctx, coolDownKey, "1", captchaCoolDown).Result()
	if err != nil {
		return errors.InternalServerError("", "set captcha coolDown: %v", err)
	}
	if !ok {
		return errors.BadRequest("", "Sending time too frequently！")
	}

	return nil

}

// 在邮箱发送失败时，清除存入的redis（验证码和限制频率发送的数据）
func (b *RegisterBIZ) cleanupCaptchaOnSendFailed(ctx context.Context, email string) {
	captchaKey := fmt.Sprintf("register:captcha:%s", email)
	cooldownKey := fmt.Sprintf("register:captcha:coolDown:%s", email)

	_ = b.redisClient.Del(ctx, captchaKey, cooldownKey).Err()
}

// 随机生成验证码
func generateRandCaptcha() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%06d", n.Int64()), nil
}

// 验证码hash加密
func hashCaptchaCode(email, code, secret string) string {
	sum := sha256.Sum256([]byte(email + ":" + code + ":" + secret))

	return hex.EncodeToString(sum[:])
}

// hash验证(不比较明文，而是重新计算hash)
func (b *RegisterBIZ) verifyHashCaptchaCode(ctx context.Context, email, code string) error {
	captchaKey := fmt.Sprintf("register:captcha:%s", email)

	savedHash, err := b.redisClient.Get(ctx, captchaKey).Result()

	if err == redis.Nil {
		return errors.BadRequest("", "Captcha Expired or not present！ ")
	}

	if err != nil {
		return errors.InternalServerError("", "get captcha: %v", err)
	}

	inputHash := hashCaptchaCode(email, code, b.captchaSecret)

	if savedHash != inputHash {
		return errors.BadRequest("", "captcha error!")
	}

	_ = b.redisClient.Del(ctx, captchaKey).Err()

	return nil

}

// 验证码存入redis
func (b *RegisterBIZ) saveCaptcha(ctx context.Context, email, code string) error {
	captchaKey := fmt.Sprintf("register:captcha:%s", email)

	hashedCode := hashCaptchaCode(email, code, b.captchaSecret)

	if err := b.redisClient.Set(ctx, captchaKey, hashedCode, captchaTTL).Err(); err != nil {
		return errors.InternalServerError("", "save captcha:%v", err)
	}

	return nil

}
