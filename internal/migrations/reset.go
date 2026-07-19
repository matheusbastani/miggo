package migrations

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

// Reset rolls back all applied migrations in reverse order.
// It executes all .down.sql files and removes all migration records.
func Reset(db *sql.DB, baseDir string) error {
	type migration struct {
		name     string
		folder   string
		downFile string
	}

	rows, err := db.Query("SELECT name FROM miggo ORDER BY applied_at DESC")
	if err != nil {
		return err
	}
	defer func() {
		_ = rows.Close()
	}()

	var appliedMigrations []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return err
		}
		appliedMigrations = append(appliedMigrations, name)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	if len(appliedMigrations) == 0 {
		color.Yellow("no migrations to reset")
		return nil
	}

	var migrations []migration

	for _, migrationName := range appliedMigrations {
		parts := strings.Split(migrationName, string(filepath.Separator))
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
			return err
		}

		sqlContent := strings.TrimSpace(string(content))
		if sqlContent == "" {
			color.Yellow("down file for %s is empty, skipping SQL execution", m.folder)
		} else {
			_, err = db.Exec(sqlContent)
			if err != nil {
				return err
			}
		}

		_, err = db.Exec("DELETE FROM miggo WHERE name = $1", m.name)
		if err != nil {
			return err
		}
	}

	color.Green("migrations reset complete")
	return nil
}

// ResetAndDrop rolls back all migrations and drops the migrations table.
func ResetAndDrop(db *sql.DB, baseDir string) error {
	Reset(db, baseDir)

	_, err := db.Exec("DROP TABLE IF EXISTS miggo")
	if err != nil {
		return err
	}

	color.Green("table miggo dropped")
	return nil
}
