package main

import (
	"log"

	"bradreed.co.uk/iverbs/api/cache"
	iverbs_http "bradreed.co.uk/iverbs/api/http"
	"bradreed.co.uk/iverbs/api/options"
	gofly "bradreed.co.uk/iverbs/gofly/gofly"
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

	server := &iverbs_http.Server{
		Port:          opts.Port,
		Fetcher:       fetcher,
		CacheProvider: cacheProvider,
	}

	log.Fatal(server.Start())
}
