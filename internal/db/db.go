package db

import "database/sql"

func Connect(name string) (*sql.DB, error) {
	conn, err := sql.Open("sqlite3", name+"?_journal=WAL&_timeout=5000&_fk=true")
	if err != nil {
		return nil, err
	}

	conn.SetMaxOpenConns(1)
	conn.SetMaxIdleConns(1)

	if err := conn.Ping(); err != nil {
		return nil, err
	}

	return conn, nil
}
