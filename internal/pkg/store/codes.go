package store

import (
	"context"
	"database/sql"
	"errors"

	cErrors "github.com/ibks-bank/profile/internal/pkg/errors"
	"github.com/ibks-bank/profile/internal/pkg/store/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (st *store) GetCode(ctx context.Context, code string) (*models.AuthenticationCode, error) {
	lastCode, err := models.AuthenticationCodes(
		models.AuthenticationCodeWhere.Code.EQ(code),
		models.AuthenticationCodeWhere.Expired.EQ(false),
		qm.OrderBy(`created_at desc`),
	).One(ctx, st.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, cErrors.Wrap(err, "can't get code")
	}

	return lastCode, nil
}

func (st *store) InsertCode(ctx context.Context, code *models.AuthenticationCode) error {
	return code.Insert(ctx, st.db, boil.Infer())
}

func (st *store) ExpireCode(ctx context.Context, code string, userID int64) error {
	lastCode, err := models.AuthenticationCodes(
		models.AuthenticationCodeWhere.Code.EQ(code),
		models.AuthenticationCodeWhere.UserID.EQ(userID),
		qm.OrderBy(`created_at desc`),
	).One(ctx, st.db)
	if err != nil {
		return cErrors.Wrap(err, "can't get code")
	}

	lastCode.Expired = true

	_, err = lastCode.Update(ctx, st.db, boil.Whitelist(models.AuthenticationCodeColumns.Expired))
	if err != nil {
		return cErrors.Wrap(err, "can't update code")
	}

	return nil
}
