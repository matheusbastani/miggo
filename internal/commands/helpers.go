package commands

import (
	"database/sql"

	"github.com/matheusbastani/miggo/internal/settings"
)

func getDatabase(name string) (*sql.DB, settings.Database, error) {
	db, set, err := settings.GetDatabase(name)
	if err != nil {
		return nil, settings.Database{}, err
	}

	return db, set, nil
}

func closeDatabase(db *sql.DB) {
	_ = db.Close()
}

func getSecure(environment settings.Environment) (bool, error) {
	return settings.IsSecure(environment)
}
