package sql

import (
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var _internalDb *sqlx.DB

func getDb() *sqlx.DB {
	if _internalDb != nil {
		return _internalDb
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	openDb, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname))
	if err != nil {
		panic(fmt.Sprintf("failed to open DB: %s", err))
	}

	if err := openDb.Ping(); err != nil {
		panic(fmt.Sprintf("failed to ping DB: %s", err))
	}

	openDb.SetMaxIdleConns(12)
	openDb.SetMaxOpenConns(20)
	openDb.SetConnMaxIdleTime(time.Hour)
	openDb.SetConnMaxLifetime(8 * time.Hour)

	fmt.Printf("Connected to database: %s:%s %s\n", host, port, dbname)
	_internalDb = openDb
	return openDb
}
