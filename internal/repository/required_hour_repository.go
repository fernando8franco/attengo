package repository

import (
	"context"
	"database/sql"

	"github.com/fernando8franco/attengo/internal/repository/sqlc"
)

type RequiredHour struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	Minutes   int    `json:"minutes"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}

type CrateRequiredHourParams struct {
	Type    string
	Minutes int
}

type RequiredHourRepository interface {
	Create(ctx context.Context, params CrateRequiredHourParams) (RequiredHour, error)
}

type requiredHourRepository struct {
	queries *sqlc.Queries
}

func NewRequiredHourRepository(db *sql.DB) RequiredHourRepository {
	return &requiredHourRepository{
		queries: sqlc.New(db),
	}
}

func (r *requiredHourRepository) Create(ctx context.Context, params CrateRequiredHourParams) (RequiredHour, error) {
	row, err := r.queries.CreateRequiredHours(
		ctx,
		sqlc.CreateRequiredHoursParams{
			Type:    params.Type,
			Minutes: int64(params.Minutes),
		},
	)
	if err != nil {
		return RequiredHour{}, err
	}

	return mapRequiredHourCreate(row), nil
}

func mapRequiredHourCreate(row sqlc.RequiredHour) RequiredHour {
	return RequiredHour{
		ID:      int(row.ID),
		Type:    row.Type,
		Minutes: int(row.Minutes),
	}
}
