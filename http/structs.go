package http

import (
  gofly "bradreed.co.uk/iverbs/gofly/gofly"
)

type LanguagesResponse struct{
  Data []*gofly.Language
  Error string `json:"error,omitempty"`
}
