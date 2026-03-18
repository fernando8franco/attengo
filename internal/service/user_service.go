package service

import (
	"context"
	"database/sql"
	"math/rand"
	"strings"
	"time"

	"github.com/fernando8franco/attengo/internal/apperr"
	"github.com/fernando8franco/attengo/internal/auth"
	"github.com/fernando8franco/attengo/internal/repository"
)

const (
	passwordLenght = 5
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
	SetUpAdmin(ctx context.Context, input CreateAdminInput) (repository.CreateAdminRow, error)
}

type userService struct {
	queries *repository.Queries
}

func NewUserService(db *sql.DB) UserService {
	return &userService{queries: repository.New(db)}
}

func (s *userService) CreateUser(ctx context.Context, input CreateUserInput) (repository.CreateUserRow, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.TrimSpace(input.Email)
	password := passwordGenetator(passwordLenght)

	row, err := s.queries.CreateUser(ctx, repository.CreateUserParams{
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

func (s *userService) SetUpAdmin(ctx context.Context, input CreateAdminInput) (repository.CreateAdminRow, error) {
	exists, err := s.queries.ExistsAdmin(ctx)
	if err != nil {
		return repository.CreateAdminRow{}, err
	}

	if exists {
		return repository.CreateAdminRow{}, apperr.NewForbiddenRequest("an admin account has already been set up")
	}

	hashedPassword, err := auth.HashPassword(input.Password)
	if err != nil {
		return repository.CreateAdminRow{}, err
	}

	row, err := s.queries.CreateAdmin(ctx, repository.CreateAdminParams{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
	})
	if err != nil {
		return repository.CreateAdminRow{}, err
	}

	return row, err
}

func passwordGenetator(length int) string {
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
