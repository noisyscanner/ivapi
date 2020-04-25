package main

import (
	iverbs_http "bradreed.co.uk/iverbs/api/http"
	gofly "bradreed.co.uk/iverbs/gofly/gofly"
	"fmt"
	"log"
)

func connect(opts *Options) (fetcher *gofly.Fetcher, err error) {
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
	opts := getOpts()

	fetcher, err := connect(opts)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}

	server := &iverbs_http.Server{
		Port:    opts.Port,
		Fetcher: fetcher,
	}

	log.Fatal(server.Start())
}
