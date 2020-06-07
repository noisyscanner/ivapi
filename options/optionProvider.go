package options

import "flag"

type Options struct {
	Port   int
	Config string
	Redis  string
}

func GetOpts() *Options {
	port := flag.Int("port", 7000, "Port for http to listen on")
	config := flag.String("config", "config", "DB config file")
	redis := flag.String("redis", "localhost:6379", "Redis host and port")

	flag.Parse()

	return &Options{
		Port:   *port,
		Config: *config,
		Redis:  *redis,
	}
}
