package miggo

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/fatih/color"
)

// Create creates a new migration directory with up and down SQL files.
// It automatically generates the next sequential index or uses the provided index.
//
// Parameters:
//   - dir: base directory where migrations are stored
//   - name: descriptive name for the migration
//   - index: optional index number (if not provided, uses next available number)
func Create(dir, name string, index ...int) {
	re := regexp.MustCompile(`^(\d{3})_`)

	entries, err := os.ReadDir(dir)
	if err != nil {
		color.Red("error reading migration directory: %s", err)
		os.Exit(1)
	}

	var indices []int
	for _, entry := range entries {
		if entry.IsDir() {
			matches := re.FindStringSubmatch(entry.Name())
			if len(matches) == 2 {
				if num, convErr := strconv.Atoi(matches[1]); convErr == nil {
					indices = append(indices, num)
				}
			}
		}
	}

	sort.Ints(indices)

	nextIndex := 1
	if len(indices) > 0 {
		nextIndex = indices[len(indices)-1] + 1
	}

	if len(index) > 0 {
		nextIndex = index[0]
	}

	prefixedName := fmt.Sprintf("%03d_%s", nextIndex, name)
	migrationDir := filepath.Join(dir, prefixedName)

	if err = os.MkdirAll(migrationDir, 0o755); err != nil {
		color.Red("error creating migration directory: %s", err)
		os.Exit(1)
	}

	timestamp := time.Now().Format("20060102150405")
	upPath := filepath.Join(migrationDir, timestamp+"_"+name+".up.sql")
	downPath := filepath.Join(migrationDir, timestamp+"_"+name+".down.sql")

	upFile, err := os.OpenFile(upPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		color.Red("error creating up migration file: %s", err)
		os.Exit(1)
	}
	defer upFile.Close()

	downFile, err := os.OpenFile(downPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		color.Red("error creating down migration file: %s", err)
		os.Exit(1)
	}
	defer downFile.Close()

	color.Green("created migration: %s", prefixedName)
}
