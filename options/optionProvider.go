package options

import (
	"os"
	"strconv"
)

type Options struct {
	Port  int
	Redis string
}

func GetOpts() *Options {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic(err)
	}

	return &Options{
		Port:  port,
		Redis: os.Getenv("REDIS"),
	}
}
