package server

import (
	"github.com/F7icK/multiplexer/internal/multiplexer/server/handlers"
	"net/http"
)

func NewRouter(h *handlers.Handlers) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/multiplexer", h.Multiplexer)

	return mux
}
