package dbHandler

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// import z _ powoduje wywołanie init() z tej package
func RunHandler(msg string) {
	db, err := sql.Open("mysql", "user7:s$cret@tcp(127.0.0.1:3306)/testdb")
	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	var version string

	err2 := db.QueryRow("SELECT VERSION()").Scan(&version)

	if err2 != nil {
		log.Fatal(err2)
	}

	fmt.Println(version)

}
