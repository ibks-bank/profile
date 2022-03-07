package profile

import (
	"context"

	"github.com/ibks-bank/profile/internal/pkg/errors"
	"github.com/ibks-bank/profile/pkg/profile"
)

func (srv *Server) SignIn(ctx context.Context, req *profile.SignInRequest) (*profile.SignInResponse, error) {
	token, err := srv.auth.SignIn(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, errors.Wrap(err, "can't sign in")
	}

	return &profile.SignInResponse{Token: token}, nil
}
