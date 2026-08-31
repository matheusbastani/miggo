package migrations

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/fatih/color"
	"github.com/matheusbastani/miggo/internal/errs"
	"github.com/matheusbastani/miggo/internal/settings"
)

// SetRollbackBoundary creates a rollback boundary on a migration.
//
// If index is empty, it uses the latest applied migration.
func SetRollbackBoundary(db *sql.DB, driver settings.Driver, index string) error {
	migration, err := findMigration(db, driver, index)
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

	_, err = db.Exec(
		fmt.Sprintf("UPDATE miggo SET rollback_boundary = TRUE WHERE migration = %s", placeholder(driver, 1)),
		migration,
	)

	if err != nil {
		return err
	}

	color.Blue("migrations locked at %s", migration)
	return nil
}

// RemoveRollbackBoundary removes the rollback boundary from a migration.
//
// If index is empty, it uses the latest applied migration.
func RemoveRollbackBoundary(db *sql.DB, driver settings.Driver, index string) error {
	if index == "" {
		return errors.New("migration must be specified")
	}

	migration, err := findMigration(db, driver, index)
	if err != nil {
		return err
	}

	_, err = db.Exec(
		fmt.Sprintf("UPDATE miggo SET rollback_boundary = FALSE WHERE migration = %s", placeholder(driver, 1)),
		migration,
	)

	color.Blue("migrations unlocked")
	return err
}

func findMigration(db *sql.DB, driver settings.Driver, index string) (string, error) {
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

	err = db.QueryRow(
		fmt.Sprintf("SELECT migration FROM miggo WHERE migration LIKE %s LIMIT 1", placeholder(driver, 1)),
		prefix+"%",
	).Scan(&migration)

	if err == sql.ErrNoRows {
		return "", errs.ErrMigrationNotFound
	}

	return migration, err
}
