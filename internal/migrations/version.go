package migrations

import (
	"database/sql"

	"github.com/fatih/color"
)

// Version displays the latest applied migration folder.
//
// If no migrations have been applied, it displays a message indicating this.
func Version(db *sql.DB) error {
	if err := createMiggoTable(db); err != nil {
		return err
	}

	var migration string
	err := db.QueryRow(`
		SELECT migration
		FROM miggo
		ORDER BY applied_at DESC
		LIMIT 1
	`).Scan(&migration)

	if err == sql.ErrNoRows {
		color.Blue("no migrations applied")
		return nil
	}

	if err != nil {
		return err
	}

	if migration == "" {
		color.Blue("no migrations applied")
		return nil
	}

	color.Blue("latest migration: %s", migration)
	return nil
}
