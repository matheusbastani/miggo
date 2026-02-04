package miggo

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
//
// Parameters:
//   - db: database connection
//   - baseDir: base directory containing migration folders
func Up(db *sql.DB, baseDir string) {
	type migration struct {
		path   string
		upFile string
		dbKey  string
	}

	var migrations []migration

	entries, err := os.ReadDir(baseDir)
	if err != nil {
		color.Red("error reading migration directory: %s", err)
		os.Exit(1)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		folderPath := filepath.Join(baseDir, entry.Name())
		files, err := os.ReadDir(folderPath)
		if err != nil {
			color.Red("error reading migration folder %s: %s", folderPath, err)
			os.Exit(1)
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
		CREATE TABLE IF NOT EXISTS schema_migrations (
			id UUID PRIMARY KEY,
			name TEXT UNIQUE NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		color.Red("error creating migrations table: %s", err)
		os.Exit(1)
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].dbKey < migrations[j].dbKey
	})

	for _, m := range migrations {
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE name = $1", m.dbKey).Scan(&count)
		if err != nil {
			color.Red("error checking for migration %s: %s", m.dbKey, err)
			os.Exit(1)
		}
		if count > 0 {
			continue
		}

		content, err := os.ReadFile(m.upFile)
		if err != nil {
			color.Red("error reading migration file %s: %s", m.upFile, err)
			os.Exit(1)
		}

		sql := strings.TrimSpace(string(content))
		if sql == "" {
			continue
		}

		tx, err := db.Begin()
		if err != nil {
			color.Red("error starting transaction for migration %s: %s", m.dbKey, err)
			os.Exit(1)
		}

		_, err = tx.Exec(sql)
		if err != nil {
			tx.Rollback()
			color.Red("error applying migration %s: %s", m.dbKey, err)
			os.Exit(1)
		}

		migrationID := uuid.New().String()
		_, err = tx.Exec("INSERT INTO schema_migrations (id, name) VALUES ($1, $2)", migrationID, m.dbKey)
		if err != nil {
			tx.Rollback()
			color.Red("error recording migration %s: %s", m.dbKey, err)
			os.Exit(1)
		}

		if err = tx.Commit(); err != nil {
			color.Red("error committing migration %s: %s", m.dbKey, err)
			os.Exit(1)
		}

		color.Green("applied migration %s", m.dbKey)
	}
}
