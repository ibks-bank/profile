package profile

import (
	"context"

	"github.com/ibks-bank/profile/pkg/profile"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srv *Server) SendCode(ctx context.Context, req *profile.SendCodeRequest) (*profile.SendCodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendCode not implemented")
}
