package http

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	gofly "github.com/noisyscanner/gofly/gofly"
)

func getLanguages(fetcher *gofly.Fetcher) httprouter.Handle {
	return jsonRoute(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
