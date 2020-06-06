package http

import (
	"log"
	"net/http"

	gofly "bradreed.co.uk/iverbs/gofly/gofly"
	"github.com/julienschmidt/httprouter"
)

func getLanguages(fetcher *gofly.Fetcher) httprouter.Handle {
	return route(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		langs, err := fetcher.GetLangs()
		if err != nil {
			// TODO: Error handling middleware
			log.Fatal(err)
			return
		}

		response := &LanguagesResponse{Data: langs}
		handleResponse(w, response)
	})
}

