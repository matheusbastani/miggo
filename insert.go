package miggo

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"

	"github.com/fatih/color"
)

// Insert creates a new migration at a specific index, renumbering existing migrations as needed.
// All migrations with index >= insertIndex will be incremented by 1.
//
// Parameters:
//   - dir: base directory where migrations are stored
//   - name: descriptive name for the migration
//   - insertIndex: the index number where the new migration should be inserted
func Insert(dir, name string, insertIndex int) {
	re := regexp.MustCompile(`^(\d{3})_`)

	entries, err := os.ReadDir(dir)
	if err != nil {
		color.Red("error reading migration directory: %s", err)
		os.Exit(1)
	}

	var folders []struct {
		index int
		name  string
	}
	for _, entry := range entries {
		if entry.IsDir() {
			matches := re.FindStringSubmatch(entry.Name())
			if len(matches) == 2 {
				if num, convErr := strconv.Atoi(matches[1]); convErr == nil {
					folders = append(folders, struct {
						index int
						name  string
					}{num, entry.Name()})
				}
			}
		}
	}

	sort.Slice(folders, func(i, j int) bool {
		return folders[i].index < folders[j].index
	})

	for i := len(folders) - 1; i >= 0; i-- {
		if folders[i].index >= insertIndex {
			oldPath := filepath.Join(dir, folders[i].name)
			newIndex := folders[i].index + 1
			newName := re.ReplaceAllString(folders[i].name, fmt.Sprintf("%03d_", newIndex))
			newPath := filepath.Join(dir, newName)

			err := os.Rename(oldPath, newPath)
			if err != nil {
				color.Red("error renaming folder %s to %s: %s", oldPath, newPath, err)
				os.Exit(1)
			}
		}
	}

	Create(dir, name, insertIndex)
}
