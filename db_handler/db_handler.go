package db_handler

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const (
	DB_FILENAME           = "logs.db"
	DB_LOG_CREATE_COMMAND = `
		CREATE TABLE IF NOT EXISTS log (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		time INTEGER NOT NULL,
		level INTEGER NOT NULL,
		ctx TEXT,
		msg TEXT
	);`
)

//Allows to store logs in DB
type DbLogger interface {
	Log(int64, int, string, string)
	LogData(*Log)
}

//Allows to read logs from DB
type DbLogReader interface {
	GetLogs() []Log
	GetLogsNewerThan(int) []Log
}

type Log struct {
	ID    string
	Time  string
	Level string
	Ctx   string
	Msg   string
}

//SQLite handler implementation
type SQLiteDb struct {
	dbHandler *sql.DB
}

func (sqlDb *SQLiteDb) GetLogsNewerThan(lastLogId int) []Log {
	rows, err := sqlDb.dbHandler.Query("SELECT * FROM log WHERE ID > ? ORDER BY id DESC", lastLogId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	data := []Log{}

	for rows.Next() {
		i := Log{}
		err = rows.Scan(&i.ID, &i.Time, &i.Level, &i.Ctx, &i.Msg)
		if err != nil {
			panic(err)
		}
		data = append(data, i)
	}
	return data
}

func (sqlDb *SQLiteDb) GetLogs() []Log {
	rows, err := sqlDb.dbHandler.Query("SELECT * FROM log ORDER BY id DESC LIMIT 500")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	data := []Log{}

	for rows.Next() {
		i := Log{}
		err = rows.Scan(&i.ID, &i.Time, &i.Level, &i.Ctx, &i.Msg)
		if err != nil {
			panic(err)
		}
		data = append(data, i)
	}
	return data
}

func (sqlDb *SQLiteDb) LogData(log *Log) {
	query, err := CreateInsertQuery(*log)
	if err != nil {
		panic(err)
	}

	stmt, err := sqlDb.dbHandler.Prepare(*query)
	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec()
	if err != nil {
		panic(err)
	}
}

func (sqlDb *SQLiteDb) Log(time int64, level int, ctx string, msg string) {
	stmt, err := sqlDb.dbHandler.Prepare("INSERT INTO log(time, level, ctx, msg) values(?,?,?,?)")
	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec(time, level, ctx, msg)
	if err != nil {
		panic(err)
	}
}

//Starts connection to SQLite and create log table if not exists
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
