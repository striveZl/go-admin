package biz

import "context"

type EmailSender interface {
	SendCaptcha(ctx context.Context, to, code string) error
}
