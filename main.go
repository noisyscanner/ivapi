package main

import (
	"log"

	"bradreed.co.uk/iverbs/api/cache"
	iverbs_http "bradreed.co.uk/iverbs/api/http"
	"bradreed.co.uk/iverbs/api/options"
	"bradreed.co.uk/iverbs/api/tokens"
	gofly "bradreed.co.uk/iverbs/gofly/gofly"
	"github.com/gomodule/redigo/redis"
)

func connect(opts *options.Options) (fetcher *gofly.Fetcher, err error) {
	configService := gofly.FileConfigService{File: opts.Config}
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
