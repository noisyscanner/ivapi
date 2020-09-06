package main

import (
	"log"
	"strconv"

	. "bradreed.co.uk/iverbs/api/helpers"
	"bradreed.co.uk/iverbs/api/options"
	"bradreed.co.uk/iverbs/api/server"
	gofly "github.com/noisyscanner/gofly/gofly"
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
