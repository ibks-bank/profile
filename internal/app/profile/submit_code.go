package profile

import (
	"context"
	"errors"

	"github.com/ibks-bank/libs/cerr"
	"github.com/ibks-bank/profile/internal/pkg/store"
	"github.com/ibks-bank/profile/pkg/profile"
	"google.golang.org/grpc/codes"
)

func (srv *Server) SubmitCode(ctx context.Context, req *profile.SubmitCodeRequest) (*profile.SubmitCodeResponse, error) {
	user, err := srv.store.GetUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, cerr.WrapMC(err, "user not found", codes.NotFound)
		}

		return nil, cerr.Wrap(err, "can't get user")

	}

	err = srv.store.ExpireCode(ctx, req.GetCode(), user.ID)
	if err != nil {
		return nil, cerr.Wrap(err, "can't expire code")
	}

	token, err := srv.auth.GetToken(req.GetEmail(), req.GetPassword(), user.HashSalt, user.ID)
	if err != nil {
		return nil, cerr.Wrap(err, "can't sign in")
	}

	return &profile.SubmitCodeResponse{Token: token}, nil
}
