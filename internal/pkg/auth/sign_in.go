package auth

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/ibks-bank/profile/internal/pkg/errors"
	"github.com/ibks-bank/profile/internal/pkg/store/models"
)

type store interface {
	GetUser(ctx context.Context, login, password string) (*models.User, error)
}

type Claims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	Password string `json:"password"`
}

type authorizer struct {
	key            string
	salt           string
	expireDuration time.Duration
	store          store
}

func NewAuthorizer(key, salt string, expireDuration time.Duration, store store) *authorizer {
	return &authorizer{key: key, salt: salt, expireDuration: expireDuration, store: store}
}

func (a *authorizer) SignIn(ctx context.Context, login, password string) (string, error) {
	password = a.HashPassword(password)

	user, err := a.store.GetUser(ctx, login, password)
	if err != nil {
		return "", errors.Wrap(err, "can't get user")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(a.expireDuration)),
			IssuedAt:  jwt.At(time.Now()),
		},
		Username: user.Email,
		Password: password,
	})

	return token.SignedString([]byte(a.key))
}

func (a *authorizer) HashPassword(password string) string {
	pwd := sha256.New()
	pwd.Write([]byte(password))
	pwd.Write([]byte(a.salt))
	return fmt.Sprintf("%x", pwd.Sum(nil))
}
