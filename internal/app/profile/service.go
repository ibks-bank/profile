package profile

import (
	"context"

	"github.com/ibks-bank/profile/internal/pb/bank_account"
	"github.com/ibks-bank/profile/internal/pkg/store/models"
	"github.com/ibks-bank/profile/pkg/profile"
)

type Server struct {
	profile.UnimplementedProfileServer

	store storeInterface
	auth  authInterface
	email emailInterface

	bankAccount bank_account.BankAccountClient
}

type storeInterface interface {
	CreateUser(ctx context.Context, user *models.User, passport *models.Passport) (int64, error)
	GetUser(ctx context.Context, login, password string) (*models.User, error)
	GetPassport(ctx context.Context, id int64) (*models.Passport, error)
	AddTelegramUsername(ctx context.Context, user *models.User, tgUsername string) error

	GetCode(ctx context.Context, code string) (*models.AuthenticationCode, error)
	InsertCode(ctx context.Context, code *models.AuthenticationCode) error
	ExpireCode(ctx context.Context, code string, userID int64) error
}

type authInterface interface {
	GetToken(login, password, salt string, userID int64) (string, error)
}

type emailInterface interface {
	Send(to, code string) error
}

func NewServer(
	store storeInterface,
	auth authInterface,
	email emailInterface,
	bankAccount bank_account.BankAccountClient,
) *Server {

	return &Server{
		store:       store,
		auth:        auth,
		email:       email,
		bankAccount: bankAccount,
	}
}
