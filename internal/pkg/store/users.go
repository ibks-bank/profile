package store

import (
	"context"
	"database/sql"
	"errors"

	cErrors "github.com/ibks-bank/profile/internal/pkg/errors"
	"github.com/ibks-bank/profile/internal/pkg/store/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func (st *store) GetUser(ctx context.Context, login, password string) (*models.User, error) {
	user, err := models.Users(
		models.UserWhere.Email.EQ(login),
		models.UserWhere.Password.EQ(password),
	).One(ctx, st.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, cErrors.Wrap(err, "can't perform select")
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
		return 0, err
	}

	return user.ID, nil
}
