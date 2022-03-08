package email

import (
	"net/smtp"

	"github.com/ibks-bank/profile/internal/pkg/errors"
)

type sender struct {
	from string
	pass string

	smtpServerHost string
	smtpServerPort string
}

func NewSender(from, pass, serverHost, serverPort string) *sender {
	return &sender{from: from, pass: pass, smtpServerHost: serverHost, smtpServerPort: serverPort}
}

func (s *sender) Send(to, code string) error {
	err := smtp.SendMail(
		s.smtpServerHost+":"+s.smtpServerPort,
		smtp.PlainAuth("", s.from, s.pass, s.smtpServerHost),
		s.from, []string{to}, []byte(code),
	)
	if err != nil {
		return errors.Wrap(err, "can't send email")
	}

	return nil
}
