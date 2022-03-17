package email

import (
	"crypto/tls"
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
	conn, err := tls.Dial("tcp", s.smtpServerHost+":"+s.smtpServerPort, &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         s.smtpServerHost,
	})
	if err != nil {
		return errors.Wrap(err, "can't dial tls")
	}

	smtpClient, err := smtp.NewClient(conn, s.smtpServerHost)
	if err != nil {
		return errors.Wrap(err, "can't create smtp client")
	}
	defer smtpClient.Quit()

	err = smtpClient.Auth(smtp.PlainAuth("", s.username, s.password, s.smtpServerHost))
	if err != nil {
		return errors.Wrap(err, "can't auth smtp")
	}

	err = smtpClient.Mail(s.username)
	if err != nil {
		return errors.Wrap(err, "can't do mail")
	}

	err = smtpClient.Rcpt(to)
	if err != nil {
		return errors.Wrap(err, "can't do rcpt")
	}

	w, err := smtpClient.Data()
	if err != nil {
		return errors.Wrap(err, "can't do data")
	}
	defer w.Close()

	_, err = w.Write([]byte(code))
	if err != nil {
		return errors.Wrap(err, "can't write code")
	}

	return nil
}
