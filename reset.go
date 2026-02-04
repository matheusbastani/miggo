package miggo

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

// Reset rolls back all applied migrations in reverse order.
// It executes all .down.sql files and removes all migration records.
//
// Parameters:
//   - db: database connection
//   - baseDir: base directory containing migration folders
func Reset(db *sql.DB, baseDir string) {
	type migration struct {
		name     string
		folder   string
		downFile string
	}

	rows, err := db.Query("SELECT name FROM migrations ORDER BY applied_at DESC")
	if err != nil {
		color.Red("error listing applied migrations: %s", err)
		os.Exit(1)
	}
	defer rows.Close()

	var appliedMigrations []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			color.Red("error scanning migration name: %s", err)
			os.Exit(1)
		}
		appliedMigrations = append(appliedMigrations, name)
	}

	if err := rows.Err(); err != nil {
		color.Red("error iterating migrations: %s", err)
		os.Exit(1)
	}

	if len(appliedMigrations) == 0 {
		color.Yellow("no migrations to reset")
		return
	}

	var migrations []migration

	for _, migrationName := range appliedMigrations {
		parts := strings.Split(migrationName, string(filepath.Separator))
		if len(parts) < 2 {
			color.Red("invalid migration name format: %s", migrationName)
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
			color.Yellow("warning: migration %s does not have a .down.sql file, skipping", folderName)
			continue
		}

		migrations = append(migrations, migration{
			name:     migrationName,
			folder:   folderName,
			downFile: downFile,
		})
	}

	for _, m := range migrations {
		content, err := os.ReadFile(m.downFile)
		if err != nil {
			color.Red("error reading down file %s: %s", m.downFile, err)
			os.Exit(1)
		}

		sqlContent := strings.TrimSpace(string(content))
		if sqlContent == "" {
			color.Yellow("down file for %s is empty, skipping SQL execution", m.folder)
		} else {
			_, err = db.Exec(sqlContent)
			if err != nil {
				color.Red("error executing down file %s: %s", m.downFile, err)
				os.Exit(1)
			}
		}

		_, err = db.Exec("DELETE FROM migrations WHERE name = $1", m.name)
		if err != nil {
			color.Red("error deleting migration %s: %s", m.name, err)
			os.Exit(1)
		}
	}

	color.Green("migrations reset complete")
}

// ResetAndDrop rolls back all migrations and drops the migrations table.
//
// Parameters:
//   - db: database connection
//   - baseDir: base directory containing migration folders
func ResetAndDrop(db *sql.DB, baseDir string) {
	Reset(db, baseDir)

	_, err := db.Exec("DROP TABLE IF EXISTS migrations")
	if err != nil {
		color.Red("error dropping migrations table: %s", err)
		os.Exit(1)
	}

	color.Green("table migrations dropped")
}
