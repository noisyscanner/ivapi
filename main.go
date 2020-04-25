package main

import (
	gofly "bradreed.co.uk/iverbs/gofly/gofly"
	iverbs_http "bradreed.co.uk/iverbs/api/http"
	"fmt"
	"log"
)


func connect() (fetcher *gofly.Fetcher, err error) {
	configService := gofly.FileConfigService{File: "config"}
	dbs := gofly.DatabaseService{ConfigService: configService}

  db, err := dbs.GetDb()
  if err != nil {
    return
  }

	fetcher = &gofly.Fetcher{Db: db}

  return
}

func main() {
  fetcher, err := connect()
  if err != nil {
    fmt.Printf("%s", err)
    return
  }

  server := &iverbs_http.Server{
    Port: 8080,
    Fetcher: fetcher,
  }

	log.Fatal(server.Start())
}
