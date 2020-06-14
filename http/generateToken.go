package http

import (
	"log"
	"net/http"

	"bradreed.co.uk/iverbs/api/tokens"
	"github.com/julienschmidt/httprouter"
)

func generateToken(tokenPersister tokens.TokenPersister) httprouter.Handle {
	return jsonRoute(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		token := tokens.GenerateToken()
		err := tokenPersister.PersistToken(token)

		var resp *TokenResponse

		if err != nil {
			resp = &TokenResponse{
				Error: "Could not persist token",
			}
		} else {
			resp = &TokenResponse{
				Token: token,
			}
		}

		json, err := resp.MarshalJSON()
		if err != nil {
			// TODO: Error handling middleware
			log.Fatal(err)
			return
		}
		w.Write(json)
		return
	})
}
