package settings

import "fmt"

type Environment string

const (
	EnvironmentDevelopment Environment = "development"
	EnvironmentProduction  Environment = "production"
)

func IsSecure(env Environment) (bool, error) {
	environment, err := normalizeEnvironment(env)
	if err != nil {
		return false, err
	}

	switch environment {
	case EnvironmentDevelopment:
		return false, nil
	case EnvironmentProduction:
		return true, nil
	default:
		return false, fmt.Errorf("unknown environment: %s", environment)
	}
}

func normalizeEnvironment(environment Environment) (Environment, error) {
	switch environment {
	case "dev", "development":
		return EnvironmentDevelopment, nil
	case "prod", "production":
		return EnvironmentProduction, nil
	default:
		return "", fmt.Errorf("unknown environment: %s", environment)
	}
}
