package store

import (
	"context"
	"database/sql"
	"strings"

	"github.com/ibks-bank/libs/auth"
	"github.com/ibks-bank/profile/internal/pkg/store/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func (st *store) GetUser(ctx context.Context, login, password string) (*models.User, error) {
	user, err := models.Users(
		models.UserWhere.Email.EQ(login),
	).One(ctx, st.db)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	if user.Password != password && user.Password != auth.HashPassword(password, user.HashSalt) {
		return nil, ErrNotFound
	}

	return user, nil
}

func (st *store) CreateUser(ctx context.Context, user *models.User, passport *models.Passport) (int64, error) {
	err := st.WithTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		err := passport.Insert(ctx, tx, boil.Infer())
		if err != nil {
			return err
		}

		user.PassportID = passport.ID

		err = user.Insert(ctx, tx, boil.Infer())
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		if strings.Contains(err.Error(), errViolatesUnique) {
			return 0, ErrAlreadyExists
		}

		return 0, err
	}

	return user.ID, nil
}
