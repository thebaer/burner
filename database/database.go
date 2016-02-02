package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type (
	DB struct {
		*sql.DB
	}
	Tx struct {
		*sql.Tx
	}
)

func Open() (*DB, error) {
	// Get database configuration
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_DB")
	dbHost := os.Getenv("DB_HOST")

	if dbUser == "" || dbPassword == "" {
		return nil, errors.New("Database user or password not set.")
	}
	if dbHost == "" {
		dbHost = "localhost"
	}
	if dbName == "" {
		dbName = "burnermail"
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=true", dbUser, dbPassword, dbHost, dbName))
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(250)

	return &DB{db}, nil
}
