package http

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	gofly "github.com/noisyscanner/gofly/gofly"
	"github.com/noisyscanner/ivapi/cache"
	"github.com/noisyscanner/ivapi/iap"
	"github.com/noisyscanner/ivapi/tokens"
)

type Server struct {
	Fetcher          *gofly.Fetcher
	CacheProvider    cache.CacheProvider
	Port             int
	TokenPersister   tokens.TokenPersister
	TokenValidator   tokens.TokenValidator
	ReceiptValidator *iap.IapValidator
	Router           *httprouter.Router
}

func (s *Server) Setup() {
	s.Router = httprouter.New()
	s.Router.PanicHandler = handlePanic
	s.Router.GET("/languages", getLanguages(s.Fetcher))
	s.Router.GET("/languages/:code", getLanguage(s))
	s.Router.POST("/tokens", generateToken(s.TokenPersister))
	s.Router.POST("/iapvalidate", iapValidate(s))
}

func (s *Server) Start() error {
	s.Setup()
	listen := fmt.Sprintf(":%d", s.Port)
	return http.ListenAndServe(listen, s.Router)
}

func response(w http.ResponseWriter, resp Response) {
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
