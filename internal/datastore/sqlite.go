package datastore

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"url-shortner/internal/logger"

	_ "github.com/mattn/go-sqlite3"
)

const expiryTime = 7 * 24 * time.Hour

type SqliteDatabase struct {
	logger   logger.Logger
	db       *sql.DB
	ctx      context.Context
	capacity int
}

type SqliteDatabaseInput struct {
	DatabaseName     string
	DatabaseCapacity int
	Logger           logger.Logger
	Ctx              context.Context
}

func NewSqliteDatabase(input *SqliteDatabaseInput) (*SqliteDatabase, error) {
	db, err := sql.Open("sqlite3", input.DatabaseName)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &SqliteDatabase{
		logger:   input.Logger,
		db:       db,
		ctx:      input.Ctx,
		capacity: input.DatabaseCapacity,
	}, nil
}

func (s *SqliteDatabase) GetURL(key string) (string, error) {
	selectQuery := "SELECT expiry_time,url FROM \"urls\" WHERE \"key\" = ? LIMIT 1"
	url := ""
	expiryTime := time.Time{}

	err := s.db.QueryRowContext(s.ctx, selectQuery, key).Scan(&expiryTime, &url)
	if errors.Is(err, sql.ErrNoRows) {
		s.logger.With("error", err).With("key", key).Error("url with key cannot be found")
		return "", err
	}

	if expiryTime.Before(time.Now()) {
		s.logger.With("key", key).Error("url has expired")
		return "", errors.New("url has expired")
	}
	return url, nil
}

func (s *SqliteDatabase) InsertURL(url string) (string, error) {
	tx, err := s.db.BeginTx(s.ctx, nil /* sql.TxOptions */)
	if err != nil {
		return "", err
	}

	key := newKey(s.capacity)
	id := newUUID()
	insertQuery := "INSERT INTO \"urls\" (id, key, url, expiry_time) VALUES (?, ?, ?, ?)"
	_, err = tx.ExecContext(s.ctx, insertQuery, id, key, url, time.Now().Add(expiryTime))
	if err != nil {
		s.logger.With("error", err).With("id", id).Error("failed to add url to database")
		return "", err
	}

	if err := tx.Commit(); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			s.logger.With("error", err).Error("could not rollback transaction")
		}
		return "", err
	}

	return key, nil
}
