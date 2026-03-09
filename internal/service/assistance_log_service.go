package service

import (
	"context"
	"database/sql"

	"github.com/fernando8franco/attengo/internal/repository"
)

type AssistanceLogInput struct {
	UserID         string
	UserPassword   string
	LogDescription string
}

type AssistanceLogService interface {
	TakeAttendance(ctx context.Context, input AssistanceLogInput) (repository.CreateEntryLogRow, error)
}

type assistanceLogService struct {
	queries *repository.Queries
}

func NewAssistanceLogService(db *sql.DB) AssistanceLogService {
	return &assistanceLogService{queries: repository.New(db)}
}

func (s *assistanceLogService) TakeAttendance(ctx context.Context, input AssistanceLogInput) (repository.CreateEntryLogRow, error) {

	return repository.CreateEntryLogRow{}, nil
}
