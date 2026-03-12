package service

import (
	"context"
	"database/sql"
	"strings"

	"github.com/fernando8franco/attengo/internal/apperr"
	"github.com/fernando8franco/attengo/internal/repository"
)

type RequiredHourInput struct {
	Type         string
	TotalMinutes int
}

type RequiredHourService interface {
	CreateRequiredHour(ctx context.Context, input RequiredHourInput) (repository.CreateRequiredHourRow, error)
}

type requiredHourService struct {
	queries *repository.Queries
}

func NewRequiredHourService(db *sql.DB) RequiredHourService {
	return &requiredHourService{queries: repository.New(db)}
}

func (s *requiredHourService) CreateRequiredHour(ctx context.Context, input RequiredHourInput) (repository.CreateRequiredHourRow, error) {
	input.Type = strings.TrimSpace(input.Type)

	row, err := s.queries.CreateRequiredHour(ctx, repository.CreateRequiredHourParams{
		Type:         input.Type,
		TotalMinutes: int64(input.TotalMinutes),
	})
	if err != nil {
		if IsUniqueConstraintError(err) {
			err = apperr.NewBadRequest(err.Error())
		}
		return repository.CreateRequiredHourRow{}, err
	}

	return row, nil
}
