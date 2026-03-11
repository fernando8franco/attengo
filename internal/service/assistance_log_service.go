package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/fernando8franco/attengo/internal/apperr"
	"github.com/fernando8franco/attengo/internal/repository"
	"github.com/google/uuid"
)

const (
	Assitance = "assistance"
)

type AssistanceLogInput struct {
	UserID       string
	UserPassword string
}

type AssistanceLogService interface {
	TakeAttendance(ctx context.Context, input AssistanceLogInput) (AttendaceDTO, error)
}

type assistanceLogService struct {
	queries *repository.Queries
}

func NewAssistanceLogService(db *sql.DB) AssistanceLogService {
	return &assistanceLogService{queries: repository.New(db)}
}

type AttendaceDTO struct {
	ID               string `json:"id"`
	EntryTime        string `json:"entry_time"`
	ExitTime         string `json:"exit_time"`
	RequiredTotal    int    `json:"required_total"`
	TotalAccumulated int    `json:"total_accumulated"`
	UserID           string `json:"user_id"`
}

func (s *assistanceLogService) TakeAttendance(ctx context.Context, input AssistanceLogInput) (AttendaceDTO, error) {
	userPsswrd, err := s.queries.ValidateUserPassword(ctx, repository.ValidateUserPasswordParams{
		ID:       input.UserID,
		Password: input.UserPassword,
	})
	if err != nil {
		return AttendaceDTO{}, err
	}

	// TODO VALIDATION OF THE CODE

	if !userPsswrd {
		return AttendaceDTO{}, apperr.NewUnauthorizedRequest("The user or password are incorrect")
	}

	lastEntry, err := s.queries.GetLastEntryLogByUser(ctx, input.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			entryLog, err := s.queries.CreateEntryLog(ctx, repository.CreateEntryLogParams{
				ID:             uuid.New().String(),
				LogDescription: Assitance,
				UserID:         input.UserID,
			})
			if err != nil {
				return AttendaceDTO{}, err
			}

			return mapEntryToAttendaceDTO(entryLog), nil
		}

		return AttendaceDTO{}, err
	}

	if time.Now().UTC().Format(time.DateOnly) != lastEntry.LogDate {
		return AttendaceDTO{}, apperr.NewBadRequest("The date is not the same")
	}

	exitLog, err := s.queries.UpdateExitLog(ctx, lastEntry.ID)
	if err != nil {
		return AttendaceDTO{}, err
	}

	return mapExitToAttendaceDTO(exitLog), nil
}

func mapEntryToAttendaceDTO(entry repository.CreateEntryLogRow) AttendaceDTO {
	return AttendaceDTO{
		ID:        entry.ID,
		EntryTime: entry.EntryTime.String,
		UserID:    entry.UserID,
	}
}

func mapExitToAttendaceDTO(exit repository.UpdateExitLogRow) AttendaceDTO {
	return AttendaceDTO{
		ID:               exit.ID,
		EntryTime:        exit.EntryTime.String,
		ExitTime:         exit.ExitTime.String,
		RequiredTotal:    int(exit.RequiredTotal),
		TotalAccumulated: int(exit.TotalAccumulated),
		UserID:           exit.UserID,
	}
}
