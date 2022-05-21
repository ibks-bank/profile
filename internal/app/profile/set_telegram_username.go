package profile

import (
	"context"
	"errors"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/ibks-bank/libs/auth"
	"github.com/ibks-bank/libs/cerr"
	"github.com/ibks-bank/profile/internal/pkg/store"
	"github.com/ibks-bank/profile/pkg/profile"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (srv *Server) SetTelegramUsername(ctx context.Context, req *profile.SetTelegramUsernameRequest) (*emptypb.Empty, error) {
	userInfo, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, cerr.WrapMC(err, "can't get user info from context", codes.Unauthenticated)
	}

	err = validateSetTelegramUsernameRequest(req)
	if err != nil {
		return nil, cerr.WrapMC(err, "validation error", codes.InvalidArgument)
	}

	user, err := srv.store.GetUser(ctx, userInfo.Username, userInfo.Password)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, cerr.WrapMC(err, "user not found", codes.NotFound)
		}

		return nil, cerr.Wrap(err, "can't get user")
	}

	err = srv.store.AddTelegramUsername(ctx, user, req.GetTgUsername())
	if err != nil {
		return nil, cerr.Wrap(err, "can't add telegram username")
	}

	return &emptypb.Empty{}, nil
}

func validateSetTelegramUsernameRequest(req *profile.SetTelegramUsernameRequest) error {
	err := validation.Validate(req, validation.NotNil)
	if err != nil {
		return err
	}

	err = validation.ValidateStruct(req,
		validation.Field(&req.TgUsername, validation.Required),
	)
	if err != nil {
		return err
	}

	return nil
}
