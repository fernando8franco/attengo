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
	ID               string  `json:"id"`
	EntryTime        string  `json:"entry_time"`
	ExitTime         string  `json:"exit_time"`
	RequiredTotal    float64 `json:"required_total"`
	TotalAccumulated float64 `json:"total_accumulated"`
	UserID           string  `json:"user_id"`
}

func (s *assistanceLogService) TakeAttendance(ctx context.Context, input AssistanceLogInput) (AttendaceDTO, error) {
	userID, err := s.queries.ValidateUserPassword(ctx, input.UserPassword)
	if err != nil {
		return AttendaceDTO{}, apperr.NewUnauthorizedRequest("The password is no valid")
	}

	lastEntry, err := s.queries.GetLastEntryLogByUser(ctx, userID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return AttendaceDTO{}, err
	}

	today := time.Now().UTC().Format(time.DateOnly)
	noEntry := errors.Is(err, sql.ErrNoRows)
	isNewDay := !noEntry && lastEntry.LogDate != today

	if noEntry || isNewDay {
		entryLog, err := s.queries.CreateEntryLog(ctx, repository.CreateEntryLogParams{
			ID:             uuid.NewString(),
			LogDescription: Assitance,
			UserID:         userID,
		})
		if err != nil {
			return AttendaceDTO{}, err
		}

		entry, err := mapEntryToAttendaceDTO(entryLog)
		if err != nil {
			return AttendaceDTO{}, err
		}
		return entry, nil
	}

	exitLog, err := s.queries.UpdateExitLog(ctx, lastEntry.ID)
	if err != nil {
		return AttendaceDTO{}, err
	}

	exit, err := mapExitToAttendaceDTO(exitLog)
	if err != nil {
		return AttendaceDTO{}, err
	}
	return exit, nil
}

func mapEntryToAttendaceDTO(entry repository.CreateEntryLogRow) (AttendaceDTO, error) {
	mexTimes, err := convertUTCToMexTime(entry.EntryTime.String)
	if err != nil {
		return AttendaceDTO{}, nil
	}

	return AttendaceDTO{
		ID:               entry.ID,
		EntryTime:        mexTimes[0],
		RequiredTotal:    minsToHours(int(entry.RequiredTotal)),
		TotalAccumulated: minsToHours(int(entry.TotalAccumulated)),
		UserID:           entry.UserID,
	}, nil
}

func mapExitToAttendaceDTO(exit repository.UpdateExitLogRow) (AttendaceDTO, error) {
	mexTimes, err := convertUTCToMexTime(exit.EntryTime.String, exit.ExitTime.String)
	if err != nil {
		return AttendaceDTO{}, nil
	}

	return AttendaceDTO{
		ID:               exit.ID,
		EntryTime:        mexTimes[0],
		ExitTime:         mexTimes[1],
		RequiredTotal:    minsToHours(int(exit.RequiredTotal)),
		TotalAccumulated: minsToHours(int(exit.TotalAccumulated)),
		UserID:           exit.UserID,
	}, nil
}

func convertUTCToMexTime(utcTimes ...string) ([]string, error) {
	mexTimes := []string{}

	mexicoZone := time.FixedZone("CST", -6*60*60)

	for _, t := range utcTimes {
		entryTimeUTC, err := time.Parse(time.TimeOnly, t)
		if err != nil {
			return mexTimes, err
		}

		mexTimes = append(mexTimes, entryTimeUTC.In(mexicoZone).Format(time.TimeOnly))
	}

	return mexTimes, nil
}
