package profile

import (
	"context"

	"github.com/ibks-bank/profile/pkg/profile"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srv *Server) Authorize(ctx context.Context, req *profile.AuthorizeRequest) (*profile.AuthorizeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Authorize not implemented")
}
