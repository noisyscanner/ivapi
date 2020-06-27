package http

import (
	gofly "github.com/noisyscanner/gofly/gofly"
)

type LanguagesResponse struct {
	Data  []*gofly.Language
	Error string `json:"error,omitempty"`
}

type TokenResponse struct {
	Token string
	Error string `json:"error,omitempty"`
}
