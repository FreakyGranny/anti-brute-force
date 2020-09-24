package storage

import (
	"context"

	_ "github.com/jackc/pgx/stdlib" //nolint
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

// SQLStorage sql storage implementation.
type SQLStorage struct {
	db *sqlx.DB
}

// New returns new sql storage.
func New(dsn string) *SQLStorage {
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("failed to load driver")
	}

	return &SQLStorage{db: db}
}

// Close close db connection.
func (s *SQLStorage) Close() error {
	return s.db.Close()
}

// AddToBlackList adds subnet to black list.
func (s *SQLStorage) AddToBlackList(ctx context.Context, ip, mask string) error {
	query := `INSERT INTO blacklist (ip, mask) VALUES($1, $2)`
	_, err := s.db.ExecContext(ctx, query, ip, mask)

	return err
}

// AddToWhiteList adds subnet to white list.
func (s *SQLStorage) AddToWhiteList(ctx context.Context, ip, mask string) error {
	query := `INSERT INTO whitelist (ip, mask) VALUES($1, $2)`
	_, err := s.db.ExecContext(ctx, query, ip, mask)

	return err
}

// RemoveFromBlackList removes subnet from black list.
func (s *SQLStorage) RemoveFromBlackList(ctx context.Context, ip, mask string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM blacklist WHERE ip = $1 AND mask = $2`, ip, mask)

	return err
}

// RemoveFromWhiteList removes subnet from white list.
func (s *SQLStorage) RemoveFromWhiteList(ctx context.Context, ip, mask string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM whitelist WHERE ip = $1 AND mask = $2`, ip, mask)

	return err
}

// GetBlackList returns list of subnet in black list.
func (s *SQLStorage) GetBlackList(ctx context.Context) ([]*IPNet, error) {
	return s.getSubnetList(ctx, "SELECT ip, mask FROM blacklist")
}

// GetWhiteList returns list of subnet in white list.
func (s *SQLStorage) GetWhiteList(ctx context.Context) ([]*IPNet, error) {
	return s.getSubnetList(ctx, "SELECT ip, mask FROM whitelist")
}

func (s *SQLStorage) getSubnetList(ctx context.Context, query string) ([]*IPNet, error) {
	res := make([]*IPNet, 0)
	rows, err := s.db.QueryxContext(ctx, query)
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var n IPNet
		if err := rows.StructScan(&n); err != nil {
			return res, err
		}
		res = append(res, &n)
	}

	return res, err
}
