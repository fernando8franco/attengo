package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/fernando8franco/attengo/internal/apperr"
	"github.com/fernando8franco/attengo/internal/repository"
)

type CreateRequiredHourInput struct {
	Type         string
	TotalMinutes int
}

type RequiredHourService interface {
	CreateRequiredHour(ctx context.Context, input CreateRequiredHourInput) (RequiredHourDTO, error)
}

type requiredHourService struct {
	queries *repository.Queries
}

func NewRequiredHourService(db *sql.DB) RequiredHourService {
	return &requiredHourService{queries: repository.New(db)}
}

func (s *requiredHourService) CreateRequiredHour(ctx context.Context, input CreateRequiredHourInput) (RequiredHourDTO, error) {
	input.Type = strings.TrimSpace(input.Type)

	row, err := s.queries.CreateRequiredHours(ctx, repository.CreateRequiredHoursParams{
		Type:         input.Type,
		TotalMinutes: int64(input.TotalMinutes),
	})
	if err != nil {
		if apperr.IsUniqueConstraint(err) {
			return RequiredHourDTO{}, &apperr.ConflictError{Resource: "required_hour"}
		}
		return RequiredHourDTO{}, fmt.Errorf("create required_hour: %w", err)
	}

	return mapRequiredHourCreate(row), nil
}

type RequiredHourDTO struct {
	ID      int    `json:"id"`
	Type    string `json:"type"`
	Minutes int    `json:"minutes"`
}

func mapRequiredHourCreate(row repository.RequiredHour) RequiredHourDTO {
	return RequiredHourDTO{
		ID:      int(row.ID),
		Type:    row.Type,
		Minutes: int(row.TotalMinutes),
	}
}
