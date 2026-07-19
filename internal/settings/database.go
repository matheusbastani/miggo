package settings

import (
	"database/sql"
	"fmt"
)

func GetDatabase(name string) (*sql.DB, string, error) {
	config, err := Get()
	if err != nil {
		return nil, "", err
	}

	dbConfig, ok := config.Databases[name]
	if !ok {
		return nil, "", fmt.Errorf("database %s not found", name)
	}

	db, err := NewDriver(
		dbConfig.Driver,
		dbConfig.URL,
	)
	if err != nil {
		return nil, "", err
	}

	return db, dbConfig.Path, nil
}
