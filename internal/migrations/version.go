package migrations

import (
	"database/sql"
	"path/filepath"

	"github.com/fatih/color"
)

// Version displays the latest applied migration folder.
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
		color.Blue("no migrations applied")
		return nil
	}

	var name string
	err = db.QueryRow("SELECT name FROM miggo ORDER BY applied_at DESC LIMIT 1").Scan(&name)
	if err == sql.ErrNoRows {
		color.Blue("no migrations applied")
		return nil
	}
	if err != nil {
		return err
	}

	if name == "" {
		color.Blue("no migrations applied")
		return nil
	}

	folderName := filepath.Base(filepath.Dir(name))
	color.Blue("latest migration folder: %s", folderName)
	return nil
}
