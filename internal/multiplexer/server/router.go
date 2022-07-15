package server

import (
	"net/http"

	"github.com/F7icK/multiplexer/internal/multiplexer/server/handlers"
)

func NewRouter(h *handlers.Handlers) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/multiplexer", h.Multiplexer)

	return mux
}
