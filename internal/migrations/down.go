package migrations

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

// Down rolls back the most recently applied migration.
//
// It executes the corresponding .down.sql file and removes the migration record.
func Down(db *sql.DB, baseDir string) error {
	var latestMigration string

	err := db.QueryRow(`
		SELECT migration
		FROM miggo
		ORDER BY applied_at DESC
		LIMIT 1
	`).Scan(&latestMigration)

	if err == sql.ErrNoRows {
		color.Yellow("No migrations to roll back")
		return nil
	}

	if err != nil {
		return err
	}

	var lockedMigration string

	err = db.QueryRow(`
		SELECT migration
		FROM miggo
		WHERE rollback_boundary = TRUE
		ORDER BY applied_at ASC
		LIMIT 1
	`).Scan(&lockedMigration)

	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if lockedMigration == latestMigration {
		return errors.New("cannot rollback below locked migration")
	}

	folderPath := filepath.Join(baseDir, latestMigration)

	files, err := os.ReadDir(folderPath)
	if err != nil {
		return err
	}

	var downFile string

	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".down.sql") {
			downFile = filepath.Join(folderPath, f.Name())
			break
		}
	}

	if downFile == "" {
		return errors.New("down migration file not found")
	}

	content, err := os.ReadFile(downFile)
	if err != nil {
		return err
	}

	sqlContent := strings.TrimSpace(string(content))

	if sqlContent != "" {
		_, err = db.Exec(sqlContent)
		if err != nil {
			return err
		}
	}

	_, err = db.Exec(
		"DELETE FROM miggo WHERE migration = $1",
		latestMigration,
	)

	if err != nil {
		return err
	}

	color.Yellow("rolled back migration: %s", latestMigration)

	return nil
}
