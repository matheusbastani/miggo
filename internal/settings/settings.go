package settings

import (
	"os"
)

type Settings struct {
	EnvFile   string              `yaml:"env_file"`
	Databases map[string]Database `yaml:"databases"`
}

type Database struct {
	Driver      Driver      `yaml:"driver"`
	URL         string      `yaml:"url"`
	Path        string      `yaml:"path"`
	Environment Environment `yaml:"environment"`
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
