package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "+_+diyor2005+_+"
	dbname   = "chatservis"
)

type DB struct {
	DB *sql.DB
}

func (db *DB) Connect() error {
	infoConnect := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	var err error
	db.DB, err = sql.Open("postgres", infoConnect)
	if err != nil {
		return err
	}
	err = db.DB.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}
