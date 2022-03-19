package auth

import (
	"context"
	"strings"

	"github.com/ibks-bank/profile/config"
	"github.com/ibks-bank/profile/internal/pkg/errors"
	"github.com/ibks-bank/profile/internal/pkg/headers"
	"github.com/ibks-bank/profile/internal/pkg/cerr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userInfo struct {
	Username string
	Password string
}

func Interceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {

	var err error

	switch info.FullMethod {
	case "/profile_pb.Profile/SignUp",
		"/profile_pb.Profile/SignIn",
		"/profile_pb.Profile/SubmitCode",
		"/profile_pb.Profile/SetAuthenticationCode":

	default:
		ctx, err = authorize(ctx)
		if err != nil {
			return nil, err
		}
	}

	h, err := handler(ctx, req)

	return h, err
}

func authorize(ctx context.Context) (context.Context, error) {
	authToken, err := headers.AuthToken(ctx)
	if err != nil {
		return ctx, status.Errorf(status.Code(err), "can't get auth token")
	}

	username, password, err := ParseToken(authToken, []byte(config.GetConfig().Auth.SigningKey))
	if err != nil {
		return ctx, status.Errorf(codes.Unauthenticated, err.Error())
	}

	return context.WithValue(ctx, headers.UserKey, userInfo{Username: username, Password: password}), nil
}

func GetUserInfo(ctx context.Context) (*userInfo, error) {
	user, ok := ctx.Value(headers.UserKey).(userInfo)
	if !ok {
		return nil, cerr.NewC("user info not found in context", codes.Unauthenticated)
	}

	return &user, nil
}
