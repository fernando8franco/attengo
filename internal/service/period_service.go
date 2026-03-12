package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/fernando8franco/attengo/internal/apperr"
	"github.com/fernando8franco/attengo/internal/repository"
)

type CreatePeriodInput struct {
	Name      string
	EntryDate string
	ExitDate  string
}

type PeriodService interface {
	CreatePeriod(ctx context.Context, input CreatePeriodInput) (repository.CreatePeriodRow, error)
}

type periodService struct {
	queries *repository.Queries
}

func NewPeriodService(db *sql.DB) PeriodService {
	return &periodService{queries: repository.New(db)}
}

func (s *periodService) CreatePeriod(ctx context.Context, input CreatePeriodInput) (repository.CreatePeriodRow, error) {
	input.Name = strings.TrimSpace(input.Name)

	entry, err1 := time.Parse(time.DateOnly, input.EntryDate)
	exit, err2 := time.Parse(time.DateOnly, input.ExitDate)
	if err1 != nil || err2 != nil {
		return repository.CreatePeriodRow{}, fmt.Errorf("Error parsing dates")
	}

	if entry.After(exit) {
		return repository.CreatePeriodRow{}, apperr.NewBadRequest("The entry time is after exit time")
	}

	row, err := s.queries.CreatePeriod(ctx, repository.CreatePeriodParams{
		Name:      input.Name,
		EntryDate: input.EntryDate,
		ExitDate:  input.ExitDate,
	})
	if err != nil {
		return repository.CreatePeriodRow{}, err
	}

	return row, nil
}
