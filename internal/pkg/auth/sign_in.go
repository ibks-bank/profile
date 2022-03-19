package auth

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
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
	expireDuration time.Duration
	store          store
}

func NewAuthorizer(key string, expireDuration time.Duration, store store) *authorizer {
	return &authorizer{key: key, expireDuration: expireDuration, store: store}
}

func (a *authorizer) GetToken(login, password, salt string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(a.expireDuration)),
			IssuedAt:  jwt.At(time.Now()),
		},
		Username: login,
		Password: HashPassword(password, salt),
	})

	return token.SignedString([]byte(a.key))
}

func HashPassword(password, salt string) string {
	pwd := sha256.New()
	pwd.Write([]byte(password))
	pwd.Write([]byte(salt))
	return fmt.Sprintf("%x", pwd.Sum(nil))
}
