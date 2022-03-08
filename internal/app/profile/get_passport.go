package profile

import (
	"context"

	"github.com/ibks-bank/profile/internal/pkg/auth"
	"github.com/ibks-bank/profile/internal/pkg/errors"
	"github.com/ibks-bank/profile/pkg/profile"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (srv *Server) GetPassport(ctx context.Context, req *emptypb.Empty) (*profile.Passport, error) {
	userInfo, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, errors.WrapMC(err, "can't get user info from context", codes.Unauthenticated)
	}

	user, err := srv.store.GetUser(ctx, userInfo.Username, userInfo.Password)
	if err != nil {
		return nil, errors.Wrap(err, "can't get user")
	}

	passport, err := srv.store.GetPassport(ctx, user.PassportID)
	if err != nil {
		return nil, errors.Wrap(err, "can't get passport")
	}

	return &profile.Passport{
		Series:     passport.Series,
		Number:     passport.Number,
		FirstName:  passport.FirstName,
		MiddleName: passport.MiddleName,
		LastName:   passport.LastName,
		IssuedBy:   passport.IssuedBy,
		IssuedAt:   timestamppb.New(passport.IssuedAt),
		Address:    passport.Address,
		Birthplace: passport.Birthplace,
		Birthdate:  timestamppb.New(passport.Birthdate),
	}, nil
}
