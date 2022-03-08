package profile

import (
	"context"

	"github.com/ibks-bank/profile/internal/pkg/errors"
	"github.com/ibks-bank/profile/internal/pkg/store/models"
	"github.com/ibks-bank/profile/pkg/profile"
)

func (srv *Server) SignUp(ctx context.Context, req *profile.SignUpRequest) (*profile.SignUpResponse, error) {
	userID, err := srv.store.CreateUser(
		ctx,
		&models.User{
			Email:    req.GetEmail(),
			Password: srv.auth.HashPassword(req.GetPassword()),
		},
		&models.Passport{
			Series:     req.GetPassport().GetSeries(),
			Number:     req.GetPassport().GetNumber(),
			FirstName:  req.GetPassport().GetFirstName(),
			MiddleName: req.GetPassport().GetMiddleName(),
			LastName:   req.GetPassport().GetLastName(),
			IssuedBy:   req.GetPassport().GetIssuedBy(),
			IssuedAt:   req.GetPassport().GetIssuedAt().AsTime(),
			Address:    req.GetPassport().GetAddress(),
			Birthplace: req.GetPassport().GetBirthplace(),
			Birthdate:  req.GetPassport().GetBirthdate().AsTime(),
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "can't create user")
	}

	return &profile.SignUpResponse{UserID: userID}, nil
}
