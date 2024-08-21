package storage

import (
	"database/sql"
	"fmt"
	"sync"
	"user/internal/user-service/enteties"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

type Database struct {
	DB *sql.DB
	m  sync.Mutex
}

func NewDB(db *sql.DB) *Database {
	return &Database{DB: db}
}

func New(connStr string) (*Database, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("openning database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("pinging database: %w", err)
	}
	return &Database{DB: db}, nil
}

func (db *Database) CreateUser(password string, email string) error {
	db.m.Lock()
	defer db.m.Unlock()

	log.Info().Msg("%s: creating user")
	var id int
	query := `INSERT INTO users(email, password, status, role) VALUES ($1, $2,$3, $4) RETURNING id`
	err := db.DB.QueryRow(query, email, password, "active", enteties.Buyer).Scan(&id)
	if err != nil {
		log.Error().Err(err).Msg("%s: unable to create user")
		return err
	}
	return nil
}
