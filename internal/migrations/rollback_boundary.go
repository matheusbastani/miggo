package migrations

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/fatih/color"
	"github.com/matheusbastani/miggo/internal/errs"
)

// SetRollbackBoundary creates a rollback boundary on a migration.
//
// If index is empty, it uses the latest applied migration.
func SetRollbackBoundary(db *sql.DB, index string) error {
	migration, err := findMigration(db, index)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		UPDATE miggo
		SET rollback_boundary = FALSE
		WHERE rollback_boundary = TRUE
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		UPDATE miggo
		SET rollback_boundary = TRUE
		WHERE migration = $1
	`, migration)

	if err != nil {
		return err
	}

	color.Yellow("migrations locked at %s", migration)
	return nil
}

// RemoveRollbackBoundary removes the rollback boundary from a migration.
//
// If index is empty, it uses the latest applied migration.
func RemoveRollbackBoundary(db *sql.DB, index string) error {
	if index == "" {
		return errors.New("migration must be specified")
	}

	migration, err := findMigration(db, index)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		UPDATE miggo
		SET rollback_boundary = FALSE
		WHERE migration = $1
	`, migration)

	color.Yellow("migrations unlocked")
	return err
}

func findMigration(db *sql.DB, index string) (string, error) {
	if index == "" {
		var migration string

		err := db.QueryRow(`
			SELECT migration
			FROM miggo
			ORDER BY applied_at DESC
			LIMIT 1
		`).Scan(&migration)

		if err == sql.ErrNoRows {
			return "", errs.ErrMigrationNotFound
		}

		return migration, err
	}

	number, err := strconv.Atoi(index)
	if err != nil {
		return "", err
	}

	prefix := fmt.Sprintf("%03d_", number)

	var migration string

	err = db.QueryRow(`
		SELECT migration
		FROM miggo
		WHERE migration LIKE $1
		LIMIT 1
	`, prefix+"%").Scan(&migration)

	if err == sql.ErrNoRows {
		return "", errs.ErrMigrationNotFound
	}

	return migration, err
}
