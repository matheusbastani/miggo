package settings

import (
	"database/sql"
	"fmt"
)

func GetDatabase(name string) (*sql.DB, Database, error) {
	config, err := Get()
	if err != nil {
		return nil, Database{}, err
	}

	settings, ok := config.Databases[name]
	if !ok {
		return nil, Database{}, fmt.Errorf("database %s not found", name)
	}

	db, err := NewDriver(
		settings.Driver,
		settings.URL,
	)
	if err != nil {
		return nil, Database{}, err
	}

	return db, settings, nil
}
