package service

import (
	"errors"

	"github.com/mattn/go-sqlite3"
)

func IsUniqueConstraintError(err error) bool {
	var sqliteErr sqlite3.Error
	if errors.As(err, &sqliteErr) {
		return sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique
	}
	return false
}
