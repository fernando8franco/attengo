package service

import (
	"context"
	"strings"

	"github.com/fernando8franco/attengo/internal/repository"
)

type CrateRequiredHourInput struct {
	Type    string
	Minutes int
}

type RequiredHourService interface {
	CreateRequiredHour(ctx context.Context, input CrateRequiredHourInput) (repository.RequiredHour, error)
}

type requiredHourService struct {
	requireHourRepo repository.RequiredHourRepository
}

func NewRequiredHourService(requiredHourRepo repository.RequiredHourRepository) RequiredHourService {
	return &requiredHourService{requireHourRepo: requiredHourRepo}
}

func (s *requiredHourService) CreateRequiredHour(ctx context.Context, input CrateRequiredHourInput) (repository.RequiredHour, error) {
	input.Type = strings.TrimSpace(input.Type)

	return s.requireHourRepo.Create(ctx, repository.CrateRequiredHourParams{
		Type:    input.Type,
		Minutes: input.Minutes,
	})
}
