package profile

import (
	"context"

	"github.com/ibks-bank/libs/cerr"
	"github.com/ibks-bank/profile/internal/pkg/store/models"
	"github.com/ibks-bank/profile/pkg/profile"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (srv *Server) SetAuthenticationCode(ctx context.Context, req *profile.SetAuthenticationCodeRequest) (*emptypb.Empty, error) {
	err := srv.store.InsertCode(ctx, &models.AuthenticationCode{UserID: req.GetUserID(), Code: req.GetCode()})
	if err != nil {
		return nil, cerr.Wrap(err, "can't insert code")
	}

	return &emptypb.Empty{}, nil
}
