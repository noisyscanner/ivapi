package server

import (
	"log"

	"bradreed.co.uk/iverbs/api/cache"
	iverbs_http "bradreed.co.uk/iverbs/api/http"
	"bradreed.co.uk/iverbs/api/options"
	"bradreed.co.uk/iverbs/api/tokens"
	"github.com/gomodule/redigo/redis"
	gofly "github.com/noisyscanner/gofly/gofly"
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

func connectToRedis(options *options.Options) (redis.Conn, error) {
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

	redisConn, err := connectToRedis(opts)
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
