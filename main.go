package main

import (
	"log"
	"os"
	"strconv"

	"bradreed.co.uk/iverbs/api/options"
	"bradreed.co.uk/iverbs/api/server"
	gofly "github.com/noisyscanner/gofly/gofly"
)

type EnvConfigService struct{}

func getEnvElse(varName string, fallback string) string {
	envVar := os.Getenv(varName)
	if envVar != "" {
		return envVar
	}

	return fallback
}

func (_ *EnvConfigService) GetConfig() *gofly.DBConfig {
	port, _ := strconv.Atoi(getEnvElse("DB_PORT", "3306"))
	return &gofly.DBConfig{
		Driver: getEnvElse("DB_DRIVER", "mysql"),
		Host:   getEnvElse("DB_HOST", "localhost"),
		User:   getEnvElse("DB_USER", "root"),
		Pass:   getEnvElse("DB_PASS", ""),
		Port:   port,
		Db:     getEnvElse("DB_NAME", "iverbs"),
	}
}

func main() {
	opts := options.GetOpts()
	goflyConfig := &EnvConfigService{}
	server := server.GetServer(opts, goflyConfig)

	log.Fatal(server.Start())

	log.Fatal(server.Start())
	log.Fatal(server.Start())
}
