package settings

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Settings struct {
	EnvFile   string              `yaml:"env_file"`
	Databases map[string]Database `yaml:"databases"`
}

type Database struct {
	Driver Driver `yaml:"driver"`
	URL    string `yaml:"url"`
	Path   string `yaml:"path"`
	Secure bool   `yaml:"secure"`
}

func CreateSettingsYAML() error {
	file, err := os.Create("miggo.yaml")
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.WriteString("databases:\n")
	if err != nil {
		return err
	}

	return nil
}

func Get() (Settings, error) {
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
