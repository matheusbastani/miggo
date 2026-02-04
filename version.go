package miggo

import (
	"database/sql"
	"os"

	"github.com/fatih/color"
)

// Version displays the latest applied migration.
// If no migrations have been applied, it displays a message indicating this.
//
// Parameters:
//   - db: database connection
func Version(db *sql.DB) {
	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.tables
			WHERE table_name = 'migrations'
		)
	`).Scan(&exists)
	if err != nil {
		color.Red("error checking for migrations table: %s", err)
		os.Exit(1)
	}

	if !exists {
		color.Blue("no migrations applied")
		return
	}

	var name string
	err = db.QueryRow("SELECT name FROM migrations ORDER BY applied_at DESC LIMIT 1").Scan(&name)
	if err == sql.ErrNoRows {
		color.Blue("no migrations applied")
		return
	}
	if err != nil {
		color.Red("error fetching latest migration: %s", err)
		os.Exit(1)
	}

	if name == "" {
		color.Blue("no migrations applied")
		return
	}

	color.Blue("latest migration: %s", name)
}
