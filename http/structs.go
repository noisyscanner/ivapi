package http

import (
	gofly "github.com/noisyscanner/gofly/gofly"
)

type Response interface {
	MarshalJSON() ([]byte, error)
}

type LanguagesResponse struct {
	Data  []*gofly.Language
	Error string `json:"error,omitempty"`
}

type TokenResponse struct {
	Token string
	Error string `json:"error,omitempty"`
}

type IapResponse struct {
	Success bool
	Error   string `json:"error,omitempty"`
}

type ReceiptBody struct {
	Receipt string `json:"receipt"`
}
