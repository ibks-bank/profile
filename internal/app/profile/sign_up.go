package profile

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/ibks-bank/profile/internal/pkg/auth"
	"github.com/ibks-bank/profile/internal/pkg/cerr"
	"github.com/ibks-bank/profile/internal/pkg/store"
	"github.com/ibks-bank/profile/internal/pkg/store/models"
	"github.com/ibks-bank/profile/pkg/profile"
	"google.golang.org/grpc/codes"
)

func (srv *Server) SignUp(ctx context.Context, req *profile.SignUpRequest) (*profile.SignUpResponse, error) {
	hashSalt := strings.ReplaceAll(uuid.New().String(), "-", "")

	userID, err := srv.store.CreateUser(
		ctx,
		&models.User{
			Email:    req.GetEmail(),
			Password: auth.HashPassword(req.GetPassword(), hashSalt),
			HashSalt: hashSalt,
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
		if errors.Is(err, store.ErrAlreadyExists) {
			return nil, cerr.WrapMC(err, "user already exists", codes.AlreadyExists)
		}

		return nil, cerr.Wrap(err, "can't create user")
	}

	return &profile.SignUpResponse{UserID: userID}, nil
}
