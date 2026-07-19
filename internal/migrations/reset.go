package migrations

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/matheusbastani/miggo/internal/errs"
)

// Reset rolls back all applied migrations in reverse order.
//
// It executes all .down.sql files and removes all migration records.
func Reset(db *sql.DB, baseDir string, secure bool, force bool) error {
	if secure {
		return errs.ErrSecureModeEnabled
	}

	if !force {
		var hasBoundary bool

		err := db.QueryRow(`
			SELECT EXISTS (
				SELECT 1
				FROM miggo
				WHERE rollback_boundary = TRUE
			)
		`).Scan(&hasBoundary)

		if err != nil {
			return err
		}

		if hasBoundary {
			return errs.ErrRollbackBoundaryExists
		}
	}

	type migration struct {
		name     string
		downFile string
	}

	rows, err := db.Query(`
		SELECT migration
		FROM miggo
		ORDER BY applied_at DESC
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

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
		folderPath := filepath.Join(baseDir, migrationName)

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
			color.Yellow(
				"warning: migration %s does not have a .down.sql file, skipping",
				migrationName,
			)
			continue
		}

		migrations = append(migrations, migration{
			name:     migrationName,
			downFile: downFile,
		})
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	for _, m := range migrations {
		content, err := os.ReadFile(m.downFile)
		if err != nil {
			_ = tx.Rollback()
			return err
		}

		sqlContent := strings.TrimSpace(string(content))

		if sqlContent != "" {
			_, err = tx.Exec(sqlContent)
			if err != nil {
				_ = tx.Rollback()
				return err
			}
		}

		_, err = tx.Exec(
			"DELETE FROM miggo WHERE migration = $1",
			m.name,
		)

		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	color.Yellow("migrations reset complete")
	return nil
}

func ResetAndDrop(
	db *sql.DB,
	baseDir string,
	secure bool,
	force bool,
) error {
	err := Reset(db, baseDir, secure, force)
	if err != nil {
		return err
	}

	_, err = db.Exec("DROP TABLE IF EXISTS miggo")
	if err != nil {
		return err
	}

	color.Yellow("table miggo dropped")

	return nil
}
