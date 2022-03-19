package profile

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/ibks-bank/profile/internal/pkg/cerr"
	"github.com/ibks-bank/profile/internal/pkg/store"
	"github.com/ibks-bank/profile/internal/pkg/store/models"
	"github.com/ibks-bank/profile/pkg/profile"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (srv *Server) SignIn(ctx context.Context, req *profile.SignInRequest) (*emptypb.Empty, error) {
	user, err := srv.store.GetUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, cerr.WrapMC(err, "user not found", codes.NotFound)
		}

		return nil, cerr.Wrap(err, "can't get user")
	}

	code := uuid.New().String()

	err = srv.email.Send(user.Email, code)
	if err != nil {
		return nil, cerr.Wrap(err, "can't send code")
	}

	err = srv.store.InsertCode(ctx, &models.AuthenticationCode{UserID: user.ID, Code: code})
	if err != nil {
		return nil, cerr.Wrap(err, "can't insert code")
	}

	return &emptypb.Empty{}, nil
}
