package db_handler

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pawel33317/coreCommunicationFramework/logger"
)

const (
	DB_FILENAME           = "logs.db"
	DB_LOG_CREATE_COMMAND = `
		CREATE TABLE IF NOT EXISTS logs (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		time INTEGER NOT NULL,
		level INTEGER NOT NULL,
		ctx TEXT,
		log TEXT
	);`
)

//Allows to store logs in DB
type DbLogger interface {
	Log(int64, logger.LogLevel, string, string)
}

//Allows to read logs from DB
type DbLogReader interface {
	GetLogs()
}

//SQLite handler implementation
type SQLiteDb struct {
	dbHandler *sql.DB
}

//Starts connection to SQLite and create logs table if not exists
func (sqlDb *SQLiteDb) Open() error {
	db, err := sql.Open("sqlite3", DB_FILENAME)
	sqlDb.dbHandler = db
	if err != nil {
		return err
	}

	if _, err := db.Exec(DB_LOG_CREATE_COMMAND); err != nil {
		return err
	}

	return nil
}

func (sqlDb *SQLiteDb) Close() {
	sqlDb.Close()
}
