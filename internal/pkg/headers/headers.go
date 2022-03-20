package headers

import (
	"context"
	"strconv"
	"strings"

	"google.golang.org/grpc/metadata"
)

const (
	UseMockKey = "X-Use-Mock"
)

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
