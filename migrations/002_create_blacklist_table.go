package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up002, Down002)
}

// Up002 up migration.
func Up002(tx *sql.Tx) error {
	_, err := tx.Exec(
		`CREATE TABLE blacklist (
			ip inet NOT NULL,
			mask inet NOT NULL,
			PRIMARY KEY (ip, mask));
	`)
	if err != nil {
		return err
	}

	return nil
}

// Down002 down migration.
func Down002(tx *sql.Tx) error {
	_, err := tx.Exec("DROP TABLE blacklist;")
	if err != nil {
		return err
	}

	return nil
}
