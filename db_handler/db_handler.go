package db_handler

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
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
	Log(int64, int, string, string)
}

//Allows to read logs from DB
type DbLogReader interface {
	GetLogs()
}

//SQLite handler implementation
type SQLiteDb struct {
	dbHandler *sql.DB
}

func (sqlDb *SQLiteDb) Log(time int64, level int, ctx string, msg string) {
	stmt, err := sqlDb.dbHandler.Prepare("INSERT INTO logs(time, level, ctx, log) values(?,?,?,?)")
	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec(time, level, ctx, msg)
	if err != nil {
		panic(err)
	}

	//TODO: remove records if more than 100
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
	sqlDb.dbHandler.Close()
}
