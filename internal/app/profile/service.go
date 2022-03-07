package profile

import (
	"context"

	"github.com/ibks-bank/profile/internal/pkg/store/models"
)

type Server struct {
	store storeInterface
	auth  authInterface
}

type storeInterface interface {
	CreateUser(ctx context.Context, user *models.User, passport *models.Passport) (int64, error)
	GetUser(ctx context.Context, login, password string) (*models.User, error)
	GetPassport(ctx context.Context, id int64) (*models.Passport, error)
}

type authInterface interface {
	HashPassword(password string) string
	SignIn(ctx context.Context, login, password string) (string, error)
}

func NewServer(store storeInterface, auth authInterface) *Server {
	return &Server{store: store, auth: auth}
}
