package migrations

import (
	"database/sql"

	"github.com/fatih/color"
)

// Version displays the latest applied migration folder.
//
// If no migrations have been applied, it displays a message indicating this.
func Version(db *sql.DB) error {
	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.tables
			WHERE table_name = 'miggo'
		)
	`).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		color.Yellow("no migrations applied")
		return nil
	}

	var migration string
	err = db.QueryRow(`
		SELECT migration
		FROM miggo
		ORDER BY applied_at DESC
		LIMIT 1
	`).Scan(&migration)

	if err == sql.ErrNoRows {
		color.Yellow("no migrations applied")
		return nil
	}

	if err != nil {
		return err
	}

	if migration == "" {
		color.Yellow("no migrations applied")
		return nil
	}

	color.Yellow("latest migration: %s", migration)
	return nil
}
