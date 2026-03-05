package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/fernando8franco/attengo/internal/repository"
)

type RequiredHourInput struct {
	Type         string
	TotalMinutes int
}

type RequiredHourService interface {
	CreateRequiredHour(ctx context.Context, input RequiredHourInput) (repository.CreateRequiredHoursRow, error)
}

type requiredHourService struct {
	queries *repository.Queries
}

func NewRequiredHourService(db *sql.DB) RequiredHourService {
	return &requiredHourService{queries: repository.New(db)}
}

func (s *requiredHourService) CreateRequiredHour(ctx context.Context, input RequiredHourInput) (repository.CreateRequiredHoursRow, error) {
	input.Type = strings.TrimSpace(input.Type)

	row, err := s.queries.CreateRequiredHours(ctx, repository.CreateRequiredHoursParams{
		Type:         input.Type,
		TotalMinutes: int64(input.TotalMinutes),
	})
	if err != nil {
		return repository.CreateRequiredHoursRow{}, fmt.Errorf("create user: %w", err)
	}

	return row, nil
}
