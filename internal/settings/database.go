package settings

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

func GetDatabase(name string) (*sql.DB, Database, error) {
	config, err := get()
	if err != nil {
		return nil, Database{}, err
	}

	settings, ok := config.Databases[name]
	if !ok {
		return nil, Database{}, fmt.Errorf("database %s not found", name)
	}

	db, err := newDriver(
		settings.Driver,
		settings.URL,
	)
	if err != nil {
		return nil, Database{}, err
	}

	return db, settings, nil
}

func get() (Settings, error) {
	file, err := os.ReadFile("miggo.yaml")
	if err != nil {
		return Settings{}, err
	}

	var settings Settings
	if err := yaml.Unmarshal(file, &settings); err != nil {
		return Settings{}, err
	}

	if len(settings.Databases) == 0 {
		return Settings{}, fmt.Errorf("no databases configured")
	}

	if settings.EnvFile != "" {
		if err := godotenv.Load(settings.EnvFile); err != nil {
			return Settings{}, fmt.Errorf("failed to load env file: %w", err)
		}
	} else {
		_ = godotenv.Load(".env")
	}

	for name, database := range settings.Databases {
		database.URL = os.ExpandEnv(database.URL)
		settings.Databases[name] = database
	}

	return settings, nil
}
