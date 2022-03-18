package profile

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	cErrors "github.com/ibks-bank/profile/internal/pkg/errors"
	"github.com/ibks-bank/profile/internal/pkg/headers"
	"github.com/ibks-bank/profile/internal/pkg/store/models"
	"github.com/ibks-bank/profile/pkg/profile"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (srv *Server) SignIn(ctx context.Context, req *profile.SignInRequest) (*emptypb.Empty, error) {
	user, err := srv.store.GetUser(ctx, req.GetEmail(), srv.auth.HashPassword(req.GetPassword()))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, cErrors.WrapMC(err, "user not found", codes.NotFound)
		}
		return nil, cErrors.Wrap(err, "can't get user")
	}

	if headers.UseMock(ctx) {
		return &emptypb.Empty{}, nil
	}

	code := uuid.New().String()

	err = srv.email.Send(user.Email, code)
	if err != nil {
		return nil, cErrors.Wrap(err, "can't send code")
	}

	err = srv.store.InsertCode(ctx, &models.AuthenticationCode{UserID: user.ID, Code: code})
	if err != nil {
		return nil, cErrors.Wrap(err, "can't insert code")
	}

	return &emptypb.Empty{}, nil
}
