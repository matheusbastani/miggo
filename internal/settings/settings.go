package settings

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Settings struct {
	Databases map[string]Database `yaml:"databases"`
}

type Database struct {
	Driver Driver `yaml:"driver"`
	URL    string `yaml:"url"`
	Path   string `yaml:"path"`
}

func CreateSettingsYAML() error {
	file, err := os.Create("miggo.yaml")
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

	_, err = file.WriteString("databases:\n")
	if err != nil {
		fmt.Println(err)
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
		return Settings{}, err
	}

	return settings, nil
}
