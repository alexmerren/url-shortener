package datastore

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"url-shortner/internal/logger"

	"time"

	"github.com/jackc/pgx/v5"
)

const (
	expiryTime                         = 7 * 24 * time.Hour
	selectExpiryTimeAndUrlWithKeyQuery = "SELECT expiry_time,url FROM \"urls\" WHERE \"key\" = ? LIMIT 1"
	insertAllIntoUrlsQuery             = "INSERT INTO \"urls\" (id, key, url, expiry_time) VALUES (?, ?, ?, ?)"
)

type PostgresUrlStore struct {
	ctx      context.Context
	connUrl  string
	conn     *pgx.Conn
	logger   logger.Logger
	capacity int
}

func NewPostgresUrlStore(
	ctx context.Context,
	logger logger.Logger,
	user, password, host, name string,
	port, capacity int,
) *PostgresUrlStore {
	if logger == nil {
		return nil
	}
	if user == "" || password == "" || host == "" || name == "" || port == 0 {
		return nil
	}

	connectionUrl := fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=disable",
		user, password, host, port, name,
	)

	return &PostgresUrlStore{
		ctx:      ctx,
		connUrl:  connectionUrl,
		logger:   logger,
		capacity: capacity,
	}
}

func (p *PostgresUrlStore) Start() error {
	p.logger.With("address", p.connUrl).Info("Starting Postgres connection")
	conn, err := pgx.Connect(p.ctx, p.connUrl)
	if err != nil {
		p.logger.WithError(err).Error("could not connect to database")
		return err
	}
	p.conn = conn
	return nil
}

func (p *PostgresUrlStore) Stop() error {
	if err := p.conn.Close(p.ctx); err != nil {
		p.logger.WithError(err).Error("could not close connection to database")
		return err
	}
	return nil
}

func (s *PostgresUrlStore) GetUrl(key string) (string, error) {
	return s.GetUrlWithContext(s.ctx, key)
}

func (s *PostgresUrlStore) GetUrlWithContext(ctx context.Context, key string) (string, error) {
	url := ""
	expiryTime := time.Time{}

	err := s.conn.QueryRow(ctx, selectExpiryTimeAndUrlWithKeyQuery, key).Scan(&expiryTime, &url)
	if errors.Is(err, sql.ErrNoRows) {
		s.logger.WithError(err).With("key", key).Error("url with key cannot be found")
		return "", err
	}

	if expiryTime.Before(time.Now()) {
		s.logger.With("key", key).Error("url has expired")
		return "", errors.New("url has expired")
	}
	return url, nil
}

func (s *PostgresUrlStore) InsertUrl(url string) (string, error) {
	return s.InsertUrlWithContext(s.ctx, url)
}

func (s *PostgresUrlStore) InsertUrlWithContext(ctx context.Context, url string) (string, error) {
	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return "", err
	}

	key := newKey(s.capacity)
	id := newUUID()
	expiryTime := time.Now().Add(expiryTime)
	_, err = tx.Exec(ctx, insertAllIntoUrlsQuery, id, key, url, expiryTime)
	if err != nil {
		s.logger.WithError(err).With("id", id).Error("failed to add url to database")
		return "", err
	}

	if err := tx.Commit(ctx); err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			s.logger.WithError(err).Error("could not rollback transaction")
		}
		return "", err
	}

	return key, nil
}
