package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up001, Down001)
}

// Up001 up migration.
func Up001(tx *sql.Tx) error {
	_, err := tx.Exec(
		`CREATE TABLE whitelist (
			ip varchar NOT NULL,
			mask varchar NOT NULL,
			PRIMARY KEY (ip, mask));
	`)
	if err != nil {
		return err
	}

	return nil
}

// Down001 down migration.
func Down001(tx *sql.Tx) error {
	_, err := tx.Exec("DROP TABLE whitelist;")
	if err != nil {
		return err
	}

	return nil
}
