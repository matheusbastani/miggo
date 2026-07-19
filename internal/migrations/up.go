package migrations

import (
	"database/sql"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/google/uuid"
)

// Up applies all pending migrations to the database.
// It creates the migrations tracking table if it doesn't exist,
// then runs all .up.sql files that haven't been applied yet.
func Up(db *sql.DB, baseDir string) error {
	type migration struct {
		path   string
		upFile string
		dbKey  string
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
					path:   folderPath,
					upFile: filepath.Join(folderPath, f.Name()),
					dbKey:  filepath.Join(entry.Name(), f.Name()),
				})
			}
		}
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS miggo (
			id UUID PRIMARY KEY,
			name TEXT UNIQUE NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].dbKey < migrations[j].dbKey
	})

	for _, m := range migrations {
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM miggo WHERE name = $1", m.dbKey).Scan(&count)
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

		migrationID := uuid.New().String()
		_, err = tx.Exec("INSERT INTO miggo (id, name) VALUES ($1, $2)", migrationID, m.dbKey)
		if err != nil {
			_ = tx.Rollback()
			return err
		}

		if err = tx.Commit(); err != nil {
			return err
		}

		color.Green("applied migration %s", m.dbKey)
	}

	return nil
}
