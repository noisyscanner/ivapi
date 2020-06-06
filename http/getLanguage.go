package http

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"bradreed.co.uk/iverbs/api/cache"
	"github.com/julienschmidt/httprouter"
)

func getLanguage(cacheProvider cache.CacheProvider) httprouter.Handle {
	return route(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		langCode := p.ByName("code")

		var (
			file *os.File
			err  error
		)
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			w.Header().Set("Content-Encoding", "gzip")
			file, err = cacheProvider.GetCacheGz(langCode)
		} else {
			file, err = cacheProvider.GetCache(langCode)
		}

		if err != nil {
			resp := &LanguagesResponse{
				Error: "Language not found",
			}
			json, err := resp.MarshalJSON()
			if err != nil {
				// TODO: Error handling middleware
				log.Fatal(err)
				return
			}
			w.Write(json)
			return
		}

		defer file.Close()

		io.Copy(w, file)
	})
}
