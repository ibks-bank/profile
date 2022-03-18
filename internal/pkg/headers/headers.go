package headers

import (
	"context"
	"strconv"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	UserKey    = "X-Auth-User"
	TokenKey   = "X-Auth-Token"
	UseMockKey = "X-Use-Mock"
)

func AuthToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Errorf(codes.InvalidArgument, "Retrieving metadata is failed")
	}

	authHeader, ok := md[strings.ToLower(TokenKey)]
	if !ok {
		return "", status.Errorf(codes.Unauthenticated, "Authorization token is not supplied")
	}

	return authHeader[0], nil
}

func UseMock(ctx context.Context) bool {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return false
	}

	useMockHeader, ok := md[strings.ToLower(UseMockKey)]
	if !ok {
		return false
	}

	useMock, err := strconv.ParseBool(useMockHeader[0])
	if err != nil {
		return false
	}

	return useMock
}
