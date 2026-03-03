package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/fernando8franco/attengo/internal/database"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./assistance.db?_journal=WAL&_timeout=5000")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dbQueries := database.New(db)

	r, err := dbQueries.CreateRequiredHours(context.Background(), database.CreateRequiredHoursParams{
		Type:    "practicas",
		Minutes: 30000,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(r)
}
