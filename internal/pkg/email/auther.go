package email

import (
	"net/smtp"

	"github.com/ibks-bank/profile/internal/pkg/errors"
)

type auther struct {
	username string
	password string
}

func loginAuth(username, password string) smtp.Auth {
	return &auther{username: username, password: password}
}

func (a *auther) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *auther) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unknown fromServer")
		}
	}
	return nil, nil
}
