package email

import (
	"net/smtp"

	"github.com/ibks-bank/profile/internal/pkg/errors"
)

type sender struct {
	username string
	password string

	smtpServerHost string
	smtpServerPort string
}

func NewSender(from, pass, serverHost, serverPort string) *sender {
	return &sender{
		username:       from,
		password:       pass,
		smtpServerHost: serverHost,
		smtpServerPort: serverPort,
	}
}

func (s *sender) Send(to, code string) error {
	auth := loginAuth(s.username, s.password)
	err := smtp.SendMail(
		s.smtpServerHost+":"+s.smtpServerPort,
		auth,
		s.username, []string{to}, []byte(code),
	)
	if err != nil {
		return errors.Wrap(err, "can't send email")
	}

	return nil
}
