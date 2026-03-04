package platform

import (
	"os"
	"strconv"
)

type config struct {
	APP_MODE     string
	ADDR         string
	ReadTimeout  int
	WriteTimeout int
	IdleTimeout  int
}

// Config automatically reads the environment variables and configures the application config
// resulting `error` can be neglacted for now
func Config() (*config, error) {
	app_mode := getEnvOrDefaultString("APP_MODE", string(APP_MODE_DEV))

	if app_mode == string(APP_MODE_DEV) {
		return &config{
			ADDR:         getEnvOrDefaultString("ADDR_DEV", "development:8080"),
			ReadTimeout:  10,
			WriteTimeout: 10,
			IdleTimeout:  10,
		}, nil
	}

	return &config{
		ADDR:         getEnvOrDefaultString("ADDR", "localhost:8080"),
		ReadTimeout:  getEnvOrDefaultInt("ReadTimeout", 10),
		WriteTimeout: getEnvOrDefaultInt("WriteTimeout", 10),
		IdleTimeout:  getEnvOrDefaultInt("IdleTimeout", 10),
	}, nil
}

func getEnvOrDefaultString(key, defVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return defVal
}

func getEnvOrDefaultInt(key string, defVal int) int {
	if val := os.Getenv(key); val != "" {
		parsed, err := strconv.Atoi(val)
		if err != nil {
			return defVal
		}
		return parsed
	}
	return defVal
}

type APP_MODE_ENUM string

const (
	APP_MODE_DEV  APP_MODE_ENUM = "APP_MODE_DEV"
	APP_MODE_PROD APP_MODE_ENUM = "APP_MODE_PROD"
)
