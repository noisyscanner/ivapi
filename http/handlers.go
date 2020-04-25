package http

import (
	gofly "bradreed.co.uk/iverbs/gofly/gofly"
	"github.com/julienschmidt/httprouter"
  "net/http"
  "fmt"
)

type Server struct{
  Fetcher *gofly.Fetcher
  Port int
}

func handleResponse(w http.ResponseWriter, response *LanguagesResponse) (err error) {
    json, err := response.MarshalJSON()
    if err != nil {
      fmt.Fprintf(w, `{"error": "Could not generate response: %s"}`, err)
      return
    }

    fmt.Fprintf(w, "%s", json)
    return
}

func handlePanic(w http.ResponseWriter, r *http.Request, err interface{}) {
  errStr, _ := err.(string)
  response := &LanguagesResponse{
    Error: errStr,
  }
  handleResponse(w, response)
}

func getLanguages(fetcher *gofly.Fetcher) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    langs, _ := fetcher.GetLangs()
    response := &LanguagesResponse{Data: langs}
    handleResponse(w, response)
  }
}


func (s *Server) Start() error {
	router := httprouter.New()
  router.PanicHandler = handlePanic
	router.GET("/languages", getLanguages(s.Fetcher))

  listen := fmt.Sprintf(":%d", s.Port)
	return http.ListenAndServe(listen, router)
}
