package profile

import (
	"context"
	"database/sql"
	"errors"

	cErrors "github.com/ibks-bank/profile/internal/pkg/errors"
	"github.com/ibks-bank/profile/pkg/profile"
	"google.golang.org/grpc/codes"
)

func (srv *Server) SubmitCode(ctx context.Context, req *profile.SubmitCodeRequest) (*profile.SubmitCodeResponse, error) {
	user, err := srv.store.GetUser(ctx, req.GetEmail(), srv.auth.HashPassword(req.GetPassword()))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, cErrors.WrapMC(err, "user not found", codes.NotFound)
		}
		return nil, cErrors.Wrap(err, "can't get user")
	}

	err = srv.store.ExpireCode(ctx, req.GetCode(), user.ID)
	if err != nil {
		return nil, cErrors.Wrap(err, "can't expire code")
	}

	token, err := srv.auth.SignIn(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, cErrors.Wrap(err, "can't sign in")
	}

	return &profile.SubmitCodeResponse{Token: token}, nil
}
