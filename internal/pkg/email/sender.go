package email

import (
	"crypto/tls"

	"github.com/ibks-bank/libs/cerr"
	gomail "gopkg.in/mail.v2"
)

type sender struct {
	username string
	password string

	smtpServerHost string
	smtpServerPort int64
}

func NewSender(from, pass, serverHost string, serverPort int64) *sender {
	return &sender{
		username:       from,
		password:       pass,
		smtpServerHost: serverHost,
		smtpServerPort: serverPort,
	}
}

func (s *sender) Send(to, code string) error {
	m := gomail.NewMessage()

	m.SetHeaders(map[string][]string{
		"From":    {s.username},
		"To":      {to},
		"Subject": {"Authentication code"},
	})
	m.SetBody("text/plain", code)

	d := gomail.NewDialer(s.smtpServerHost, int(s.smtpServerPort), s.username, s.password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	err := d.DialAndSend(m)
	if err != nil {
		return cerr.Wrap(err, "can't dial and send")
	}

	return nil
}
