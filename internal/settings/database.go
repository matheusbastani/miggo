package settings

import (
	"database/sql"
	"fmt"
	"os"
	"reflect"

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
		if err := loadEnv(); err != nil {
			return Settings{}, err
		}
	}

	expandEnv(&settings)

	return settings, nil
}

func loadEnv() error {
	files := []string{
		".env.local",
		".env",
	}

	for _, file := range files {
		if _, err := os.Stat(file); err == nil {
			if err := godotenv.Load(file); err != nil {
				return err
			}
		}
	}

	return nil
}

func expandEnv(value any) {
	v := reflect.ValueOf(value)

	if v.Kind() != reflect.Ptr || v.IsNil() {
		return
	}

	expandValue(v.Elem())
}

func expandValue(v reflect.Value) {
	switch v.Kind() {
	case reflect.String:
		if v.CanSet() {
			v.SetString(os.ExpandEnv(v.String()))
		}

	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			expandValue(v.Field(i))
		}

	case reflect.Map:
		for _, key := range v.MapKeys() {
			value := v.MapIndex(key)

			copy := reflect.New(value.Type()).Elem()
			copy.Set(value)

			expandValue(copy)

			v.SetMapIndex(key, copy)
		}

	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			expandValue(v.Index(i))
		}
	}
}
