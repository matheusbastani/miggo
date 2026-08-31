package migrations

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/matheusbastani/miggo/internal/settings"
)

// Up applies all pending migrations to the database.
//
// It creates the migrations tracking table if it doesn't exist.
func Up(db *sql.DB, driver settings.Driver, baseDir string) error {
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
			fmt.Sprintf("SELECT COUNT(*) FROM miggo WHERE migration = %s", placeholder(driver, 1)),
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

		sqlContent := strings.TrimSpace(string(content))
		if sqlContent == "" {
			continue
		}

		tx, err := db.Begin()
		if err != nil {
			return err
		}

		_, err = tx.Exec(sqlContent)
		if err != nil {
			_ = tx.Rollback()
			return err
		}

		_, err = tx.Exec(
			fmt.Sprintf("INSERT INTO miggo (migration) VALUES (%s)", placeholder(driver, 1)),
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

// placeholder returns the parameter placeholder syntax for the given driver.
// Postgres uses $1, $2, ...; SQLite and MySQL use ?.
func placeholder(driver settings.Driver, n int) string {
	if driver == settings.DriverPostgres {
		return fmt.Sprintf("$%d", n)
	}
	return "?"
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

	return nil
}
