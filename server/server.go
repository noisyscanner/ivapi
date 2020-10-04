package server

import (
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
	gofly "github.com/noisyscanner/gofly/gofly"
	"github.com/noisyscanner/ivapi/cache"
	iverbs_http "github.com/noisyscanner/ivapi/http"
	"github.com/noisyscanner/ivapi/options"
	"github.com/noisyscanner/ivapi/tokens"
)

func connect(configService gofly.ConfigService) (fetcher *gofly.Fetcher, err error) {
	dbs := gofly.DatabaseService{ConfigService: configService}

	db, err := dbs.GetDb()
	if err != nil {
		return
	}

	fetcher = &gofly.Fetcher{Db: db}

	return
}

func ConnectToRedis(options *options.Options) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", options.Redis)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func GetServer(opts *options.Options, goflyConfig gofly.ConfigService) (server *iverbs_http.Server) {
	fetcher, err := connect(goflyConfig)
	if err != nil {
		log.Fatal(err)
		return
	}

	cacheProvider := &cache.FileCacheProvider{
		RootDirectory: opts.CacheDirectory,
	}

	if err != nil {
		log.Fatal(err)
		return
	}

	redisPool := ConnectToRedis(opts)

	tokenPersister := tokens.NewRedisTokenPersister(redisPool)
	tokenValidator := tokens.NewRedisTokenValidator(redisPool)

	return &iverbs_http.Server{
		Port:           opts.Port,
		Fetcher:        fetcher,
		CacheProvider:  cacheProvider,
		TokenPersister: tokenPersister,
		TokenValidator: tokenValidator,
	}
}
