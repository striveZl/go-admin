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
	"math/big"
	"time"

	"github.com/go-redis/redis/v8"
)

type EmailSender interface {
	SendCaptcha(ctx context.Context, to, code string) error
}

type CaptchaBIZ struct {
	userDAl       *dal.UserDAL
	redisClient   *redis.Client
	emailSender   EmailSender
	captchaSecret string
}

func NewCaptcha(userDAl *dal.UserDAL, redisClient *redis.Client, emailSender EmailSender,
	captchaSecret string) *CaptchaBIZ {
	return &CaptchaBIZ{
		userDAl:       userDAl,
		redisClient:   redisClient,
		emailSender:   emailSender,
		captchaSecret: captchaSecret,
	}
}

// 频率
const captchaCoolDown = time.Second * 60

// 验证码过期时间
const captchaTTL = time.Minute * 5

// 验证码发送
func (captchaBIZ *CaptchaBIZ) Captcha(ctx context.Context, req *schema.CaptchaRequest) (*schema.CaptchaResponse, error) {
	if captchaBIZ == nil || captchaBIZ.userDAl == nil || captchaBIZ.redisClient == nil || captchaBIZ.emailSender == nil {
		return nil, errors.InternalServerError("", "Captcha dal is not initialized")
	}

	existsByEmail, err := captchaBIZ.userDAl.ExistsByEmail(ctx, req.Email)

	if err != nil {
		return nil, errors.InternalServerError("", "Email:%v", err)
	}

	if existsByEmail {
		return nil, errors.BadRequest("", "The Email already register")
	}

	if err := captchaBIZ.checkCaptchaCoolDown(ctx, req.Email); err != nil {
		return nil, err
	}

	code, err := generateRandCaptcha()

	if err != nil {
		return nil, errors.InternalServerError("", "generate captcha: %v", err)
	}

	if err := captchaBIZ.saveCaptcha(ctx, req.Email, code); err != nil {
		return nil, err
	}

	if err := captchaBIZ.emailSender.SendCaptcha(ctx, req.Email, code); err != nil {
		captchaBIZ.cleanupCaptchaOnSendFailed(ctx, req.Email)
		return nil, errors.InternalServerError("", "send captcha email: %v", err)
	}

	return &schema.CaptchaResponse{
		ExpireSeconds: int(captchaTTL.Seconds()),
	}, nil

}

// 限制发送频率
func (b *CaptchaBIZ) checkCaptchaCoolDown(ctx context.Context, email string) error {
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
func (b *CaptchaBIZ) cleanupCaptchaOnSendFailed(ctx context.Context, email string) {
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
func (b *CaptchaBIZ) verifyHashCaptchaCode(ctx context.Context, email, code string) error {
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
func (b *CaptchaBIZ) saveCaptcha(ctx context.Context, email, code string) error {
	captchaKey := fmt.Sprintf("register:captcha:%s", email)

	hashedCode := hashCaptchaCode(email, code, b.captchaSecret)

	if err := b.redisClient.Set(ctx, captchaKey, hashedCode, captchaTTL).Err(); err != nil {
		return errors.InternalServerError("", "save captcha:%v", err)
	}

	return nil

}
