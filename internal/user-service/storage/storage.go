package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"user/internal/user-service/enteties"

	"github.com/lib/pq"
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

var (
	ErrEmailOrLoginAlreadyExists = errors.New("email or login already exists")
	ErrUserNotFound              = errors.New("user not found")
)

func (db *Database) CreateUser(password, email, login string) error {
	db.m.Lock()
	defer db.m.Unlock()

	log.Info().Msg("%s: creating user")
	var id int
	query := `INSERT INTO users(email, password, login, status, role) VALUES ($1, $2,$3, $4, $5) RETURNING id`
	err := db.DB.QueryRow(query, email, password, login, "active", enteties.Buyer).Scan(&id)
	if err != nil {
		pgErr, ok := err.(*pq.Error)
		if ok {
			if pgErr.Code == "23505" {
				return ErrEmailOrLoginAlreadyExists
			}
		}
		log.Error().Err(err).Msg("%s: unable to create user")
		return err
	}
	return nil
}

func (db *Database) GetPasswordByEmail(email string) (string, error) {
	query := `SELECT password  
			 FROM users 
			 WHERE email=$1`
	var password string
	err := db.DB.QueryRow(query, email).Scan(&password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrUserNotFound
		}
		return "", err
	}
	return password, nil
}

func (db *Database) GetUserByEmail(email string) (enteties.UserPersonalInfo, error) {
	query := `SELECT name, lastname, email, login   
			 FROM users 
			 WHERE email=$1`

	var name, lastname, emailDB, login sql.NullString
	err := db.DB.QueryRow(query, email).Scan(&name, &lastname, &emailDB, &login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return enteties.UserPersonalInfo{}, ErrUserNotFound
		}
		if err != nil {
			log.Debug().Msg(err.Error())
		}
		return enteties.UserPersonalInfo{}, err
	}
	user := enteties.UserPersonalInfo{
		Name:     safeGetStringFromNull(name),
		LastName: safeGetStringFromNull(lastname),
		Email:    safeGetStringFromNull(emailDB),
		Login:    safeGetStringFromNull(login),
	}

	return user, nil
}

func safeGetStringFromNull(s sql.NullString) string {
	if s.Valid && s.String != "" {
		return s.String
	}
	return ""
}
