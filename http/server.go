package http

import (
	"fmt"
	"net/http"

	"bradreed.co.uk/iverbs/api/cache"
	gofly "bradreed.co.uk/iverbs/gofly/gofly"
	"github.com/julienschmidt/httprouter"
)

type Server struct {
	Fetcher       *gofly.Fetcher
	CacheProvider cache.CacheProvider
	Port          int
}

func (s *Server) Start() error {
	router := httprouter.New()
	router.PanicHandler = handlePanic
	router.GET("/languages", getLanguages(s.Fetcher))
	router.GET("/languages/:code", getLanguage(s.CacheProvider))

	listen := fmt.Sprintf(":%d", s.Port)
	return http.ListenAndServe(listen, router)
}
