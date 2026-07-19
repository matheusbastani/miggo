package migrations

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"

	"github.com/matheusbastani/miggo/internal/errs"
)

// Insert creates a new migration at a specific index, renumbering existing migrations as needed.
//
// All migrations with index >= insertIndex will be incremented by 1.
func Insert(
	db *sql.DB,
	dir string,
	migration string,
	insertIndex int,
	secure bool,
	force bool,
) error {
	if secure {
		return errs.ErrSecureModeEnabled
	}

	if !force {
		var locked bool

		err := db.QueryRow(`
			SELECT EXISTS (
				SELECT 1
				FROM miggo
				WHERE rollback_boundary = TRUE
			)
		`).Scan(&locked)

		if err != nil {
			return err
		}

		if locked {
			return errs.ErrRollbackBoundaryExists
		}
	}

	re := regexp.MustCompile(`^(\d{3})_`)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	var folders []struct {
		index int
		name  string
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		matches := re.FindStringSubmatch(entry.Name())
		if len(matches) != 2 {
			continue
		}

		num, err := strconv.Atoi(matches[1])
		if err != nil {
			continue
		}

		folders = append(
			folders,
			struct {
				index int
				name  string
			}{
				index: num,
				name:  entry.Name(),
			},
		)
	}

	sort.Slice(folders, func(i, j int) bool {
		return folders[i].index < folders[j].index
	})

	for i := len(folders) - 1; i >= 0; i-- {
		if folders[i].index >= insertIndex {
			oldPath := filepath.Join(dir, folders[i].name)

			newIndex := folders[i].index + 1

			newName := re.ReplaceAllString(
				folders[i].name,
				fmt.Sprintf("%03d_", newIndex),
			)

			newPath := filepath.Join(dir, newName)

			err := os.Rename(oldPath, newPath)
			if err != nil {
				return err
			}
		}
	}

	return Create(dir, migration, insertIndex)
}
