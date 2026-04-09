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
	UsersString        = "users"
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

type LoginAdminInput struct {
	Email    string
	Password string
}

type UserService interface {
	CreateUser(ctx context.Context, input CreateUserInput) (NotAdminInfo, error)
	SetUpAdmin(ctx context.Context, input CreateAdminInput) (TokensReponse, error)
	AdminLogin(ctx context.Context, input LoginAdminInput) (TokensReponse, error)
	AdminLogout(ctx context.Context, token string) error
	GetActiveUsers(ctx context.Context) ([]repository.GetActiveUsersRow, error)
	GetHoursPeriodsAndUsers(ctx context.Context) ([]repository.GetRequiredHoursRow, []repository.GetPeriodsRow, []NotAdminInfo, error)
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

type NotAdminInfo struct {
	ID         string
	Name       string
	Email      string
	Period     string
	Type       string
	Hours      int64
	TotalHours float64
	Password   string
}

func (s *userService) CreateUser(ctx context.Context, input CreateUserInput) (NotAdminInfo, error) {
	passwords, err := s.queries.GetUsersPasswords(ctx)
	if err != nil {
		return NotAdminInfo{}, err
	}

	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.TrimSpace(input.Email)
	password := passwordGenerator(PasswordLength)
	attemps := 0
	for slices.Contains(passwords, password) {
		password = passwordGenerator(PasswordLength)
		attemps++
		if attemps == MaxPasswordAttemps {
			return NotAdminInfo{}, errors.New("Max passwords attemps")
		}
	}

	user, err := s.queries.CreateUser(ctx, repository.CreateUserParams{
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
		return NotAdminInfo{}, err
	}

	hours := int64(minsToHours(int(user.Hours)))
	totalHours := minsToHours(int(user.TotalHours))
	return NotAdminInfo{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		Period:     user.Period,
		Type:       user.Type,
		Hours:      hours,
		TotalHours: totalHours,
		Password:   user.Password,
	}, nil
}

type TokensReponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *userService) SetUpAdmin(ctx context.Context, input CreateAdminInput) (TokensReponse, error) {
	exists, err := s.queries.ExistsAdmin(ctx)
	if err != nil {
		return TokensReponse{}, err
	}

	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.TrimSpace(input.Email)

	if exists {
		return TokensReponse{}, apperr.NewForbiddenRequest("an admin account has already been set up")
	}

	hashedPassword, err := auth.HashPassword(input.Password)
	if err != nil {
		return TokensReponse{}, err
	}

	adminID, err := s.queries.CreateAdmin(ctx, repository.CreateAdminParams{
		ID:       uuid.NewString(),
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
	})
	if err != nil {
		return TokensReponse{}, err
	}

	accessToken, err := auth.MakeJWT(s.cfg.IssuerJWT, s.cfg.SecretJWT, adminID, s.cfg.ExpirationTime)
	if err != nil {
		return TokensReponse{}, err
	}

	refreshToken := auth.MakeRefreshToken()

	err = s.queries.CreateRefreshToken(ctx, repository.CreateRefreshTokenParams{
		Token:  refreshToken,
		UserID: adminID,
	})
	if err != nil {
		return TokensReponse{}, err
	}

	response := TokensReponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return response, err
}

func (s *userService) AdminLogin(ctx context.Context, input LoginAdminInput) (TokensReponse, error) {
	user, err := s.queries.GetAdminIDAndPasswordByEmail(ctx, input.Email)
	if err != nil {
		return TokensReponse{}, apperr.NewBadRequest(err.Error())
	}

	match, err := auth.CheckPasswordHash(input.Password, user.Password)
	if err != nil || !match {
		return TokensReponse{}, apperr.NewUnauthorizedRequest("Incorrect email or password")
	}

	accessToken, err := auth.MakeJWT(s.cfg.IssuerJWT, s.cfg.SecretJWT, user.ID, s.cfg.ExpirationTime)
	if err != nil {
		return TokensReponse{}, err
	}

	refreshToken := auth.MakeRefreshToken()

	err = s.queries.CreateRefreshToken(ctx, repository.CreateRefreshTokenParams{
		Token:  refreshToken,
		UserID: user.ID,
	})

	reponse := TokensReponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return reponse, nil
}

func (s *userService) AdminLogout(ctx context.Context, token string) error {
	err := s.queries.SetRevokedAt(ctx, token)
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) GetActiveUsers(ctx context.Context) ([]repository.GetActiveUsersRow, error) {
	users, err := s.queries.GetActiveUsers(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *userService) GetHoursPeriodsAndUsers(ctx context.Context) ([]repository.GetRequiredHoursRow, []repository.GetPeriodsRow, []NotAdminInfo, error) {
	hours, err := s.queries.GetRequiredHours(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	periods, err := s.queries.GetPeriods(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	users, err := s.queries.GetNotAdminUsers(ctx)
	if err != nil {
		return nil, nil, nil, err
	}
	usersDTO := []NotAdminInfo{}

	for _, user := range users {
		hours := int64(minsToHours(int(user.Hours)))
		totalHours := minsToHours(int(user.TotalHours))

		usersDTO = append(usersDTO, NotAdminInfo{
			ID:         user.ID,
			Name:       user.Name,
			Email:      user.Email,
			Period:     user.Period,
			Type:       user.Type,
			Hours:      hours,
			TotalHours: totalHours,
			Password:   user.Password,
		})
	}

	return hours, periods, usersDTO, nil
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
