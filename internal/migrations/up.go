package migrations

import (
	"database/sql"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fatih/color"
)

// Up applies all pending migrations to the database.
//
// It creates the migrations tracking table if it doesn't exist.
func Up(db *sql.DB, baseDir string) error {
	err := createMiggoTable(db)
	if err != nil {
		return err
	}

	type migration struct {
		name   string
		upFile string
	}

	var migrations []migration

	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		folderPath := filepath.Join(baseDir, entry.Name())

		files, err := os.ReadDir(folderPath)
		if err != nil {
			return err
		}

		for _, f := range files {
			if !f.IsDir() && strings.HasSuffix(f.Name(), ".up.sql") {
				migrations = append(migrations, migration{
					name:   entry.Name(),
					upFile: filepath.Join(folderPath, f.Name()),
				})
			}
		}
	}

	if len(migrations) == 0 {
		color.Blue("database is up to date")
		return nil
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].name < migrations[j].name
	})

	applied := false

	for _, m := range migrations {
		var count int
		err := db.QueryRow(
			"SELECT COUNT(*) FROM miggo WHERE migration = $1",
			m.name,
		).Scan(&count)

		if err != nil {
			return err
		}

		if count > 0 {
			continue
		}

		content, err := os.ReadFile(m.upFile)
		if err != nil {
			return err
		}

		sql := strings.TrimSpace(string(content))
		if sql == "" {
			continue
		}

		tx, err := db.Begin()
		if err != nil {
			return err
		}

		_, err = tx.Exec(sql)
		if err != nil {
			_ = tx.Rollback()
			return err
		}

		_, err = tx.Exec(
			"INSERT INTO miggo (migration) VALUES ($1)",
			m.name,
		)

		if err != nil {
			_ = tx.Rollback()
			return err
		}

		if err = tx.Commit(); err != nil {
			return err
		}

		applied = true
		color.Blue("applied migration %s", m.name)
	}

	if !applied {
		color.Blue("database is up to date")
	}

	return nil
}

func createMiggoTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS miggo (
			migration TEXT PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			rollback_boundary BOOLEAN DEFAULT FALSE
		)
	`)

	if err != nil {
		return err
	}

	color.Blue("miggo table created")
	return nil
}
