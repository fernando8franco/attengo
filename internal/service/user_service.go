package service

import (
	"context"
	"database/sql"
	"errors"
	"math/rand"
	"slices"
	"strings"
	"time"

	"github.com/fernando8franco/attengo/internal/apperr"
	"github.com/fernando8franco/attengo/internal/auth"
	"github.com/fernando8franco/attengo/internal/config"
	"github.com/fernando8franco/attengo/internal/repository"
	"github.com/google/uuid"
)

const (
	PasswordLength     = 6
	MaxPasswordAttemps = 10000
)

type CreateUserInput struct {
	Name           string
	Email          string
	RequiredHourID int
	PeriodID       int
}

type CreateAdminInput struct {
	Name     string
	Email    string
	Password string
}

type UserService interface {
	CreateUser(ctx context.Context, input CreateUserInput) (repository.CreateUserRow, error)
	SetUpAdmin(ctx context.Context, input CreateAdminInput) (SetUpAdminReponse, error)
}

type userService struct {
	queries *repository.Queries
	cfg     *config.Config
}

func NewUserService(db *sql.DB, cfg *config.Config) UserService {
	return &userService{
		queries: repository.New(db),
		cfg:     cfg,
	}
}

func (s *userService) CreateUser(ctx context.Context, input CreateUserInput) (repository.CreateUserRow, error) {
	passwords, err := s.queries.GetUsersPasswords(ctx)
	if err != nil {
		return repository.CreateUserRow{}, err
	}

	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.TrimSpace(input.Email)
	password := passwordGenerator(PasswordLength)
	attemps := 0
	for slices.Contains(passwords, password) {
		password = passwordGenerator(PasswordLength)
		attemps++
		if attemps == MaxPasswordAttemps {
			return repository.CreateUserRow{}, errors.New("Max passwords attemps")
		}
	}

	row, err := s.queries.CreateUser(ctx, repository.CreateUserParams{
		ID:       uuid.NewString(),
		Name:     input.Name,
		Email:    input.Email,
		Password: password,
		RequiredHourID: sql.NullInt64{
			Int64: int64(input.RequiredHourID),
			Valid: true,
		},
		PeriodID: sql.NullInt64{
			Int64: int64(input.PeriodID),
			Valid: true,
		},
	})
	if err != nil {
		if IsUniqueConstraintError(err) {
			err = apperr.NewBadRequest(err.Error())
		}
		return repository.CreateUserRow{}, err
	}

	return row, nil
}

type SetUpAdminReponse struct {
	repository.CreateAdminRow
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *userService) SetUpAdmin(ctx context.Context, input CreateAdminInput) (SetUpAdminReponse, error) {
	exists, err := s.queries.ExistsAdmin(ctx)
	if err != nil {
		return SetUpAdminReponse{}, err
	}

	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.TrimSpace(input.Email)

	if exists {
		return SetUpAdminReponse{}, apperr.NewForbiddenRequest("an admin account has already been set up")
	}

	hashedPassword, err := auth.HashPassword(input.Password)
	if err != nil {
		return SetUpAdminReponse{}, err
	}

	admin, err := s.queries.CreateAdmin(ctx, repository.CreateAdminParams{
		ID:       uuid.NewString(),
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
	})
	if err != nil {
		return SetUpAdminReponse{}, err
	}

	accessToken, err := auth.MakeJWT(s.cfg.IssuerJWT, s.cfg.SecretJWT, admin.ID, s.cfg.ExpirationTime)
	if err != nil {
		return SetUpAdminReponse{}, err
	}

	refreshToken := auth.MakeRefreshToken()

	_, err = s.queries.CreateRefreshToken(ctx, repository.CreateRefreshTokenParams{
		Token:  refreshToken,
		UserID: admin.ID,
	})
	if err != nil {
		return SetUpAdminReponse{}, err
	}

	return SetUpAdminReponse{
		CreateAdminRow: admin,
		AccessToken:    accessToken,
		RefreshToken:   refreshToken,
	}, err
}

func passwordGenerator(length int) string {
	lowerCase := "abcdefghijklmnopqrstuvwxyz"
	number := "0123456789"

	var password strings.Builder

	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	for range length {
		randNum := rng.Intn(2)

		switch randNum {
		case 0:
			randCharNum := rng.Intn(len(lowerCase))
			password.WriteString(string(lowerCase[randCharNum]))
		case 1:
			randCharNum := rng.Intn(len(number))
			password.WriteString(string(number[randCharNum]))
		}
	}

	return password.String()
}
