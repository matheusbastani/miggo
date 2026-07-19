package migrations

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
func Create(dir, name string, index ...int) error {
	re := regexp.MustCompile(`^(\d{3})_`)

	if dir == "" {
		return fmt.Errorf("migration directory not specified")
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
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
		return err
	}

	timestamp := time.Now().Format("20060102150405")

	upPath := filepath.Join(migrationDir, timestamp+"_"+name+".up.sql")
	downPath := filepath.Join(migrationDir, timestamp+"_"+name+".down.sql")

	upFile, err := os.OpenFile(upPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer upFile.Close()

	downFile, err := os.OpenFile(downPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer downFile.Close()

	color.Green("created migration: %s", prefixedName)

	return nil
}
