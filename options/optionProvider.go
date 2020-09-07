package options

import (
	"os"
	"strconv"

	"github.com/noisyscanner/ivapi/helpers"
)

type Options struct {
	Port           int
	Redis          string
	CacheDirectory string
}

func GetOpts() *Options {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic(err)
	}

	return &Options{
		Port:           port,
		Redis:          os.Getenv("REDIS"),
		CacheDirectory: helpers.GetEnvElse("CACHE_DIRECTORY", "../gofly/langcache"),
	}
}
