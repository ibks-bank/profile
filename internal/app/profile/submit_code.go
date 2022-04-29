package profile

import (
	"context"
	"errors"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/ibks-bank/libs/cerr"
	"github.com/ibks-bank/profile/internal/pkg/store"
	"github.com/ibks-bank/profile/pkg/profile"
	"google.golang.org/grpc/codes"
)

func (srv *Server) SubmitCode(ctx context.Context, req *profile.SubmitCodeRequest) (*profile.SubmitCodeResponse, error) {
	err := validateSubmitCodeRequest(req)
	if err != nil {
		return nil, cerr.WrapMC(err, "validation error", codes.InvalidArgument)
	}

	user, err := srv.store.GetUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, cerr.WrapMC(err, "user not found", codes.NotFound)
		}

		return nil, cerr.Wrap(err, "can't get user")

	}

	code, err := srv.store.GetCode(ctx, req.GetCode())
	if err != nil && !errors.Is(err, store.ErrNotFound) {
		return nil, cerr.Wrap(err, "can't get code")
	}

	if code.Expired {
		return nil, cerr.NewC("code already used", codes.InvalidArgument)
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

func validateSubmitCodeRequest(req *profile.SubmitCodeRequest) error {
	err := validation.Validate(req, validation.NotNil)
	if err != nil {
		return err
	}

	err = validation.ValidateStruct(req,
		validation.Field(&req.Email, validation.Required),
		validation.Field(&req.Password, validation.Required),
		validation.Field(&req.Code, validation.Required),
	)
	if err != nil {
		return err
	}

	return nil
}
