package main

import "flag"

type Options struct {
	Port   int
	Config string
}

func getOpts() *Options {
	port := flag.Int("port", 7000, "Port for http to listen on")
	config := flag.String("config", "config", "DB config file")

	flag.Parse()

	return &Options{
		Port:   *port,
		Config: *config,
	}
}
