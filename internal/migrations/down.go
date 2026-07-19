package migrations

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

// Down rolls back the most recently applied migration.
// It executes the corresponding .down.sql file and removes the migration record.
func Down(db *sql.DB, baseDir string) error {
	var latestMigration string
	err := db.QueryRow("SELECT name FROM miggo ORDER BY applied_at DESC LIMIT 1").Scan(&latestMigration)
	if err == sql.ErrNoRows {
		return err
	}
	if err != nil {
		return err
	}

	if latestMigration == "" {
		color.Yellow("No migrations to roll back")
		return nil
	}

	parts := strings.Split(latestMigration, string(filepath.Separator))
	if len(parts) < 2 {
		return err
	}

	folderName := parts[0]
	folderPath := filepath.Join(baseDir, folderName)

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
		return err
	}

	content, err := os.ReadFile(downFile)
	if err != nil {
		return err
	}

	sqlContent := strings.TrimSpace(string(content))
	if sqlContent == "" {
		color.Yellow("Down file for %s is empty, skipping SQL execution", folderName)
	} else {
		_, err = db.Exec(sqlContent)
		if err != nil {
			return err
		}
	}

	_, err = db.Exec("DELETE FROM miggo WHERE name = $1", latestMigration)
	if err != nil {
		return err
	}

	color.Green("rolled back migration: %s", folderName)
	return nil
}
