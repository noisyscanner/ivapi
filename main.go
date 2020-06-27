package main

import (
	"log"
	"os"
	"strconv"

	"bradreed.co.uk/iverbs/api/cache"
	iverbs_http "bradreed.co.uk/iverbs/api/http"
	"bradreed.co.uk/iverbs/api/options"
	"bradreed.co.uk/iverbs/api/tokens"
	"github.com/gomodule/redigo/redis"
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

func connect(opts *options.Options) (fetcher *gofly.Fetcher, err error) {
	configService := &EnvConfigService{}
	dbs := gofly.DatabaseService{ConfigService: configService}

	db, err := dbs.GetDb()
	if err != nil {
		return
	}

	fetcher = &gofly.Fetcher{Db: db}

	return
}

func connectToRedis(options *options.Options) (redis.Conn, error) {
	return redis.Dial("tcp", options.Redis)
}

func main() {
	opts := options.GetOpts()

	fetcher, err := connect(opts)
	if err != nil {
		log.Fatal(err)
		return
	}

	cacheProvider := &cache.FileCacheProvider{
		RootDirectory: "../gofly/langcache",
	}

	// err, redisTokenValidator := tokens.NewRedisTokenValidator()
	if err != nil {
		log.Fatal(err)
		return
	}

	redisConn, err := connectToRedis(opts)
	if err != nil {
		log.Fatal(err)
		return
	}

	tokenPersister := tokens.NewRedisTokenPersister(redisConn)
	tokenValidator := tokens.NewRedisTokenValidator(redisConn)

	server := &iverbs_http.Server{
		Port:           opts.Port,
		Fetcher:        fetcher,
		CacheProvider:  cacheProvider,
		TokenPersister: tokenPersister,
		TokenValidator: tokenValidator,
	}

	log.Fatal(server.Start())
}
