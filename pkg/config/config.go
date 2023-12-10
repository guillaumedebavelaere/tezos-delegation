package config

import (
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	// ViperEnvironment is the yaml root property that defines current application environment.
	ViperEnvironment = "environment"
	// EnvDev is the development environment, automatically loading environment files `.env` in projects root.
	EnvDev = "dev"
)

const path = "config"

// Parse parses a configuration.
func Parse(name string, config interface{}) error {
	v := viper.NewWithOptions()

	v.SetConfigName(name)
	appNameUpper := strings.ToUpper(name)

	v.SetConfigType("yaml")

	folders := []string{getEnv(appNameUpper+"_CONFIG_PATH", path)}
	for _, folderPath := range folders {
		v.AddConfigPath(folderPath)
	}

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// env prefix will be upper-cased automatically
	v.SetEnvPrefix(appNameUpper)
	v.AutomaticEnv()

	// Read the configuration file
	if err := v.ReadInConfig(); err != nil {
		zap.L().Error("Error reading config file", zap.Error(err))

		return err
	}

	configureDotEnv(v)

	// Unmarshal the configuration into a struct
	if err := v.Unmarshal(config); err != nil {
		zap.L().Error("Error unmarshalling config", zap.Error(err))

		return err
	}

	return nil
}

// Validate validates a configuration.
func Validate(config interface{}) error {
	v := validator.New()

	return v.Struct(config)
}

func configureDotEnv(v *viper.Viper) {
	if v.Get(ViperEnvironment) == EnvDev {
		if err := godotenv.Load(); err != nil {
			zap.L().Info("no .env file available on working directory, skipping...")
		} else {
			zap.L().Info(".env file successfully loaded")
		}
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
