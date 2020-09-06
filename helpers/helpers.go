package helpers

import "os"

func GetEnvElse(varName string, fallback string) string {
	envVar := os.Getenv(varName)
	if envVar != "" {
		return envVar
	}

	return fallback
}
