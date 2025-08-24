package storage

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func (storage *Storage) Open() error {

	connectionString := os.Getenv("DATABASE_URL")

	if connectionString == "" {
		log.Fatal("DATABASE_URL is empty")
	}

	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	storage.db = db
	log.Println("database connection success")
	return nil

}

func (storage *Storage) Close() {
	storage.db.Close()
}

func (s *Storage) GetDB() *sql.DB {
	return s.db
}
