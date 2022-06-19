package db_handler

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func RunDb() (int, error) {

	const create string = `
  CREATE TABLE IF NOT EXISTS activities (
  id INTEGER NOT NULL PRIMARY KEY,
  time DATETIME NOT NULL,
  description TEXT
  );`

	const file string = "activities.db"
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return 0, err
	}
	if _, err := db.Exec(create); err != nil {
		return 0, err
	}

	return 1, nil
}
