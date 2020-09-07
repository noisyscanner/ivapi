package server

import (
	"log"

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

func ConnectToRedis(options *options.Options) (redis.Conn, error) {
	return redis.Dial("tcp", options.Redis)
}

func GetServer(opts *options.Options, goflyConfig gofly.ConfigService) (server *iverbs_http.Server) {
	fetcher, err := connect(goflyConfig)
	if err != nil {
		log.Fatal(err)
		return
	}

	cacheProvider := &cache.FileCacheProvider{
		RootDirectory: "../gofly/langcache",
	}

	if err != nil {
		log.Fatal(err)
		return
	}

	redisConn, err := ConnectToRedis(opts)
	if err != nil {
		log.Fatal(err)
		return
	}

	tokenPersister := tokens.NewRedisTokenPersister(redisConn)
	tokenValidator := tokens.NewRedisTokenValidator(redisConn)

	return &iverbs_http.Server{
		Port:           opts.Port,
		Fetcher:        fetcher,
		CacheProvider:  cacheProvider,
		TokenPersister: tokenPersister,
		TokenValidator: tokenValidator,
	}
}
