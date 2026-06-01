package email

import (
	"context"
	"fmt"

	"github.com/resend/resend-go/v3"
)

type SMTPSender struct {
	EmailApiKey string
	From        string
}

func NewSMTPSender(emailApiKey, from string) *SMTPSender {
	return &SMTPSender{
		EmailApiKey: emailApiKey,
		From:        from,
	}
}

func (s *SMTPSender) SendCaptcha(ctx context.Context, to, code string) error {
	// addr := fmt.Sprintf("%s:%d", s.Host, s.Port)

	// auth := smtp.PlainAuth("", s.UserName, s.Password, s.Host)

	// subject := "Subject:注册验证码\r\n"

	// contextType := "Content-Type:text/plain; charset=UTF-8\r\n\r\n"

	// body := fmt.Sprintf("您注册的验证码是：%s，五分钟内有效。", code)

	// msg := []byte(subject + contextType + body)

	// fmt.Printf("验证信息：addr:%s------auth:%s-------From:%s-----To:%s-----Msg:%s", addr, auth, s.From, []string{to}, string(msg))

	// return smtp.SendMail(addr, auth, s.From, []string{to}, msg)

	client := resend.NewClient(s.EmailApiKey)

	body := fmt.Sprintf("<p>您注册的验证码是：<strong>%s</strong>，五分钟内有效。</p>", code)

	params := &resend.SendEmailRequest{
		From:    s.From,
		To:      []string{to},
		Subject: "注册验证码",
		Html:    body,
	}

	_, err := client.Emails.Send(params)
	if err != nil {
		return err
	}
	return nil

}
