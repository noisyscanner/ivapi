package cache

import (
	"fmt"
	"os"
	"path"
)

type CacheProvider interface {
	GetCache(langCache string) (*os.File, error)
	GetCacheGz(langCache string) (*os.File, error)
}

type FileCacheProvider struct {
	RootDirectory string
}

func (c *FileCacheProvider) GetCache(langCache string) (*os.File, error) {
	filename := fmt.Sprintf("%s.json.full", langCache)
	filePath := path.Join(c.RootDirectory, filename)
	return os.Open(filePath)
}

func (c *FileCacheProvider) GetCacheGz(langCache string) (*os.File, error) {
	filename := fmt.Sprintf("%s.json.full.gz", langCache)
	filePath := path.Join(c.RootDirectory, filename)
	return os.Open(filePath)
}
