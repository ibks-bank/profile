package auth

import (
	"context"
	"github.com/ibks-bank/profile/config"
	"github.com/ibks-bank/profile/internal/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

const (
	UserKey  = "X-Auth-User"
	TokenKey = "X-Auth-Token"
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
	case "/profile_pb.Profile/SignUp", "/profile_pb.Profile/SignIn", "/profile_pb.Profile/SubmitCode":
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
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, status.Errorf(codes.InvalidArgument, "Retrieving metadata is failed")
	}

	authHeader, ok := md[strings.ToLower(TokenKey)]
	if !ok {
		return ctx, status.Errorf(codes.Unauthenticated, "Authorization token is not supplied")
	}

	token := authHeader[0]

	username, password, err := ParseToken(token, []byte(config.GetConfig().Auth.SigningKey))
	if err != nil {
		return ctx, status.Errorf(codes.Unauthenticated, err.Error())
	}

	return context.WithValue(ctx, UserKey, userInfo{Username: username, Password: password}), nil
}

func GetUserInfo(ctx context.Context) (*userInfo, error) {
	user, ok := ctx.Value(UserKey).(userInfo)
	if !ok {
		return nil, errors.NewC("user info not found in context", codes.Unauthenticated)
	}

	return &user, nil
}
