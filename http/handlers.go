package http

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

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

func route(handler httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		header := w.Header()
		header.Set("Content-type", "application/json")
		handler(w, r, p)
	}
}

