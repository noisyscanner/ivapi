package http

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func iapValidate(s *Server) httprouter.Handle {
	return jsonRoute(s.authMiddleware(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var (
			receipt []byte
			success bool
			err     error
			errStr  string
		)
		_, err = r.Body.Read(receipt)

		if err == nil {
			success, err = s.ReceiptValidator.ValidateIapToken(receipt)
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
