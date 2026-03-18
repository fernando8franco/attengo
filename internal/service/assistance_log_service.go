package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/fernando8franco/attengo/internal/apperr"
	"github.com/fernando8franco/attengo/internal/repository"
)

const (
	Assitance = "assistance"
)

type AssistanceLogInput struct {
	UserID       int
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
	ID               int    `json:"id"`
	EntryTime        string `json:"entry_time"`
	ExitTime         string `json:"exit_time"`
	RequiredTotal    int    `json:"required_total"`
	TotalAccumulated int    `json:"total_accumulated"`
	UserID           int    `json:"user_id"`
}

func (s *assistanceLogService) TakeAttendance(ctx context.Context, input AssistanceLogInput) (AttendaceDTO, error) {
	userPsswrd, err := s.queries.ValidateUserPassword(ctx, repository.ValidateUserPasswordParams{
		ID:       int64(input.UserID),
		Password: input.UserPassword,
	})
	if err != nil {
		return AttendaceDTO{}, err
	}

	if !userPsswrd {
		return AttendaceDTO{}, apperr.NewUnauthorizedRequest("The user or password are incorrect")
	}

	lastEntry, err := s.queries.GetLastEntryLogByUser(ctx, int64(input.UserID))
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return AttendaceDTO{}, err
	}

	today := time.Now().UTC().Format(time.DateOnly)
	noEntry := errors.Is(err, sql.ErrNoRows)
	isNewDay := !noEntry && lastEntry.LogDate != today

	if noEntry || isNewDay {
		entryLog, err := s.queries.CreateEntryLog(ctx, repository.CreateEntryLogParams{
			LogDescription: Assitance,
			UserID:         int64(input.UserID),
		})
		if err != nil {
			return AttendaceDTO{}, err
		}

		return mapEntryToAttendaceDTO(entryLog), nil
	}

	exitLog, err := s.queries.UpdateExitLog(ctx, lastEntry.ID)
	if err != nil {
		return AttendaceDTO{}, err
	}

	return mapExitToAttendaceDTO(exitLog), nil
}

func mapEntryToAttendaceDTO(entry repository.CreateEntryLogRow) AttendaceDTO {
	return AttendaceDTO{
		ID:               int(entry.ID),
		EntryTime:        entry.EntryTime.String,
		RequiredTotal:    int(entry.RequiredTotal),
		TotalAccumulated: int(entry.TotalAccumulated),
		UserID:           int(entry.UserID),
	}
}

func mapExitToAttendaceDTO(exit repository.UpdateExitLogRow) AttendaceDTO {
	return AttendaceDTO{
		ID:               int(exit.ID),
		EntryTime:        exit.EntryTime.String,
		ExitTime:         exit.ExitTime.String,
		RequiredTotal:    int(exit.RequiredTotal),
		TotalAccumulated: int(exit.TotalAccumulated),
		UserID:           int(exit.UserID),
	}
}
