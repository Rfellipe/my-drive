package lib

import (
	"database/sql"
	"fmt"
	"os"

	 _ "github.com/lib/pq"
)

const (
	host = "localhost"
	port = 5432
)

type Database struct {
	Connection *sql.DB
}

var DB Database

func StartDatabaseConnection() Database {
	psqlconn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", 
		host, port, os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		panic(err)
	}

	return Database{Connection: db}
}
