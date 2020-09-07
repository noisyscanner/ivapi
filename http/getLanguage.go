package http

import (
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func getLanguage(s *Server) httprouter.Handle {
	return jsonRoute(s.authMiddleware(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		langCode := p.ByName("code")

		var (
			file *os.File
			err  error
		)
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			w.Header().Set("Content-Encoding", "gzip")
			file, err = s.CacheProvider.GetCacheGz(langCode)
		} else {
			file, err = s.CacheProvider.GetCache(langCode)
		}

		if err != nil {
			response(w, &LanguagesResponse{
				Error: "Language not found",
			})
			return
		}

		defer file.Close()

		io.Copy(w, file)
	}))
}
