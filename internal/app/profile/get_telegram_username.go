package profile

import (
	"context"
	"errors"

	"github.com/ibks-bank/libs/auth"
	"github.com/ibks-bank/libs/cerr"
	"github.com/ibks-bank/profile/internal/pkg/store"
	"github.com/ibks-bank/profile/pkg/profile"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (srv *Server) GetTelegramUsername(ctx context.Context, _ *emptypb.Empty) (*profile.GetTelegramUsernameResponse, error) {
	userInfo, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, cerr.WrapMC(err, "can't get user info from context", codes.Unauthenticated)
	}

	user, err := srv.store.GetUser(ctx, userInfo.Username, userInfo.Password)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, cerr.WrapMC(err, "user not found", codes.NotFound)
		}

		return nil, cerr.Wrap(err, "can't get user")
	}

	return &profile.GetTelegramUsernameResponse{TgUsername: user.TGUsername.String}, nil
}
