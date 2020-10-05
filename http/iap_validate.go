package http

import (
	json "encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func iapValidate(s *Server) httprouter.Handle {
	return jsonRoute(s.authMiddleware(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var (
			body        []byte
			receiptBody *ReceiptBody
			success     bool
			err         error
			errStr      string
		)
		receiptBody = &ReceiptBody{}

		body, err = ioutil.ReadAll(r.Body)

		if err == nil {
			err = json.Unmarshal(body, receiptBody)
		}

		if err == nil {
			success, err = s.ReceiptValidator.ValidateIapToken(receiptBody.Receipt)
		}

		if err != nil {
			errStr = err.Error()
		}

		response(w, &IapResponse{
			Success: success,
			Error:   errStr,
		})
	}))
}
