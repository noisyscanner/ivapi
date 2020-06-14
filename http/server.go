package http

import (
	"fmt"
	"log"
	"net/http"

	"bradreed.co.uk/iverbs/api/cache"
	"bradreed.co.uk/iverbs/api/tokens"
	gofly "bradreed.co.uk/iverbs/gofly/gofly"
	"github.com/julienschmidt/httprouter"
)

type Server struct {
	Fetcher        *gofly.Fetcher
	CacheProvider  cache.CacheProvider
	Port           int
	TokenPersister tokens.TokenPersister
	TokenValidator tokens.TokenValidator
}

func (s *Server) Start() error {
	router := httprouter.New()
	router.PanicHandler = handlePanic
	router.GET("/languages", getLanguages(s.Fetcher))
	router.GET("/languages/:code", getLanguage(s))
	router.POST("/tokens", generateToken(s.TokenPersister))

	listen := fmt.Sprintf(":%d", s.Port)
	return http.ListenAndServe(listen, router)
}

func response(w http.ResponseWriter, resp *LanguagesResponse) {
	json, err := resp.MarshalJSON()
	if err != nil {
		// TODO: Error handling middleware
		log.Fatal(err)
		return
	}
	w.Write(json)
}

func (s *Server) authMiddleware(handler httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		token := r.Header.Get("Authorization")
		if token == "" {
			response(w, &LanguagesResponse{
				Error: "Token required",
			})
			return
		}

		isValid, err := s.TokenValidator.Validate(token)
		if err != nil {
			log.Printf("Could not validate token: %v", err)
			response(w, &LanguagesResponse{
				Error: "Error validating token",
			})
			return
		}

		if !isValid {
			response(w, &LanguagesResponse{
				Error: "Invalid token",
			})
			return
		}

		err = s.TokenPersister.InvalidateToken(token)
		if err != nil {
			log.Panicf("Could not invalidate token: %v", err)
		}

		handler(w, r, p)
	}
}
