package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/aleroxac/docker-healthcheck/internal/entity"
)

type Config struct {
	Protocol string
	Host     string
	Port     int
	Path     string
}

func Init() (*Config, error) {
	needed_envs := reflect.VisibleFields(reflect.TypeOf(entity.HealthcheckConfigs{}))

	missing_envs := 0
	for _, env := range needed_envs {
		env_var_name := fmt.Sprintf("HEALTHCHECK_%v", strings.ToUpper(env.Name))
		if _, present := os.LookupEnv(env_var_name); !present {
			fmt.Printf("Please provide the %s environment variable\n", env_var_name)
			missing_envs = missing_envs + 1
		}
	}

	if missing_envs != 0 {
		return nil, nil
	}

	cfg := entity.NewHealthcheckConfigs(
		os.Getenv("HEALTHCHECK_PROTOCOL"),
		os.Getenv("HEALTHCHECK_HOST"),
		os.Getenv("HEALTHCHECK_PORT"),
		os.Getenv("HEALTHCHECK_PATH"),
	)

	config := Config{
		Protocol: cfg.Protocol,
		Host:     cfg.Host,
		Port:     cfg.Port,
		Path:     cfg.Path,
	}

	return &config, nil
}
