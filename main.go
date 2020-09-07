package main

import (
	"log"
	"strconv"

	gofly "github.com/noisyscanner/gofly/gofly"
	. "github.com/noisyscanner/ivapi/helpers"
	"github.com/noisyscanner/ivapi/options"
	"github.com/noisyscanner/ivapi/server"
)

type EnvConfigService struct{}

func (_ *EnvConfigService) GetConfig() *gofly.DBConfig {
	port, _ := strconv.Atoi(GetEnvElse("DB_PORT", "3306"))
	return &gofly.DBConfig{
		Driver: GetEnvElse("DB_DRIVER", "mysql"),
		Host:   GetEnvElse("DB_HOST", "localhost"),
		User:   GetEnvElse("DB_USER", "root"),
		Pass:   GetEnvElse("DB_PASS", ""),
		Port:   port,
		Db:     GetEnvElse("DB_NAME", "iverbs"),
	}
}

func main() {
	opts := options.GetOpts()
	goflyConfig := &EnvConfigService{}
	server := server.GetServer(opts, goflyConfig)

	log.Fatal(server.Start())
}
