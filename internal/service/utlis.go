package service

import (
	"errors"
	"math"

	"github.com/mattn/go-sqlite3"
)

func IsUniqueConstraintError(err error) bool {
	var sqliteErr sqlite3.Error
	if errors.As(err, &sqliteErr) {
		return sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique
	}
	return false
}

func minsToHours(mins int) float64 {
	return math.Round((float64(mins)/60.0)*100) / 100
}
