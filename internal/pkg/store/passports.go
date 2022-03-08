package store

import (
	"context"
	"database/sql"
	"github.com/ibks-bank/profile/internal/pkg/errors"

	"github.com/ibks-bank/profile/internal/pkg/store/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func (st *store) CreatePassport(ctx context.Context, passport *models.Passport) (int64, error) {
	err := passport.Insert(ctx, st.db, boil.Infer())
	if err != nil {
		return 0, err
	}

	return passport.ID, nil
}

func (st *store) GetPassport(ctx context.Context, id int64) (*models.Passport, error) {
	passport, err := models.FindPassport(ctx, st.db, id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, errors.Wrap(ErrNotFound, "passport not found in database")
		default:
			return nil, errors.Wrap(err, "can't find application")
		}
	}

	return passport, nil
}
