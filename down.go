package miggo

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

// Down rolls back the most recently applied migration.
// It executes the corresponding .down.sql file and removes the migration record.
//
// Parameters:
//   - db: database connection
//   - baseDir: base directory containing migration folders
func Down(db *sql.DB, baseDir string) {
	var latestMigration string
	err := db.QueryRow("SELECT name FROM schema_migrations ORDER BY applied_at DESC LIMIT 1").Scan(&latestMigration)
	if err == sql.ErrNoRows {
		color.Yellow("no migrations to roll back")
		return
	}
	if err != nil {
		color.Red("error getting latest migration: %s", err)
		os.Exit(1)
	}

	if latestMigration == "" {
		color.Yellow("No migrations to roll back")
		return
	}

	parts := strings.Split(latestMigration, string(filepath.Separator))
	if len(parts) < 2 {
		color.Red("invalid migration name format: %s", latestMigration)
		os.Exit(1)
	}

	folderName := parts[0]
	folderPath := filepath.Join(baseDir, folderName)

	files, err := os.ReadDir(folderPath)
	if err != nil {
		color.Red("error reading migration folder %s: %s", folderPath, err)
		os.Exit(1)
	}

	var downFile string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".down.sql") {
			downFile = filepath.Join(folderPath, f.Name())
			break
		}
	}

	if downFile == "" {
		color.Red("migration %s does not have a .down.sql file", folderName)
		os.Exit(1)
	}

	content, err := os.ReadFile(downFile)
	if err != nil {
		color.Red("error reading down file %s: %s", downFile, err)
		os.Exit(1)
	}

	sqlContent := strings.TrimSpace(string(content))
	if sqlContent == "" {
		color.Yellow("Down file for %s is empty, skipping SQL execution", folderName)
	} else {
		_, err = db.Exec(sqlContent)
		if err != nil {
			color.Red("error executing down file %s: %s", downFile, err)
			os.Exit(1)
		}
	}

	_, err = db.Exec("DELETE FROM schema_migrations WHERE name = $1", latestMigration)
	if err != nil {
		color.Red("error deleting migration %s: %s", latestMigration, err)
		os.Exit(1)
	}

	color.Green("rolled back migration: %s", folderName)
}
