package http

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type IapHandler func(http.ResponseWriter, *http.Request, httprouter.Params) *IapResponse

func iapHandlerWrapper(handler IapHandler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		result := handler(w, r, p)
		response(w, result)
	}
}

func iapValidate(s *Server) httprouter.Handle {
	return jsonRoute(iapHandlerWrapper(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) (res *IapResponse) {
		res = &IapResponse{}
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			log.Print("Could not read iapvalidate body", err)
			res.Error = err.Error()
			return
		}

		receiptBody := &ReceiptBody{}
		err = json.Unmarshal(body, receiptBody)

		if err != nil {
			log.Print("Could not parse iapvalidate body", err)
			res.Error = err.Error()
			return
		}

		success, err := s.ReceiptValidator.ValidateIapToken([]byte(receiptBody.Receipt))

		if err != nil {
			res.Error = err.Error()
		}

		res.Success = success
		return
	}))
}
