package profile

import (
	"context"

	"github.com/ibks-bank/profile/internal/pkg/store/models"
)

type Server struct {
	store storeInterface
	auth  authInterface
	email emailInterface
}

type storeInterface interface {
	CreateUser(ctx context.Context, user *models.User, passport *models.Passport) (int64, error)
	GetUser(ctx context.Context, login, password string) (*models.User, error)
	GetPassport(ctx context.Context, id int64) (*models.Passport, error)

	GetCode(ctx context.Context, code string) (*models.AuthenticationCode, error)
	InsertCode(ctx context.Context, code *models.AuthenticationCode) error
	ExpireCode(ctx context.Context, code string, userID int64) error
}

type authInterface interface {
	HashPassword(password string) string
	SignIn(ctx context.Context, login, password string) (string, error)
}

type emailInterface interface {
	Send(to, code string) error
}

func NewServer(store storeInterface, auth authInterface, email emailInterface) *Server {
	return &Server{store: store, auth: auth, email: email}
}
