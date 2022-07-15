package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"syscall"

	"github.com/F7icK/multiplexer/internal/multiplexer/service"
	"github.com/F7icK/multiplexer/internal/multiplexer/types"
	"github.com/F7icK/multiplexer/pkg/infrastruct"
)

type Handlers struct {
	service *service.Service
}

func NewHandlers(srv *service.Service) *Handlers {
	return &Handlers{
		service: srv,
	}
}

type result struct {
	Err string `json:"error"`
}

func errorEncode(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if customError, ok := err.(*infrastruct.CustomError); ok {
		w.WriteHeader(customError.Code)
	}

	r := result{Err: err.Error()}

	if err = json.NewEncoder(w).Encode(r); err != nil && !errors.Is(err, syscall.EPIPE) {
		log.Printf("err with Encode in errorEncode: %s", err)
	}
}

func responseEncoder(w http.ResponseWriter, res interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)

	if err := json.NewEncoder(w).Encode(res); err != nil && !errors.Is(err, syscall.EPIPE) {
		log.Printf("err with Encode in responseEncoder: %s", err)
	}
}

func (h *Handlers) Multiplexer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorEncode(w, infrastruct.ErrMethodNotAllowed)
		return
	}

	arrURLs := types.URLsRequest{URLs: nil}
	if err := json.NewDecoder(r.Body).Decode(&arrURLs); err != nil {
		errorEncode(w, infrastruct.ErrBadRequest)
		return
	}

	countURL := len(arrURLs.URLs)
	if countURL == 0 || countURL > 20 {
		errorEncode(w, infrastruct.ErrCountURL)
		return
	}

	for _, urlInUrls := range arrURLs.URLs {
		_, err := url.ParseRequestURI(urlInUrls)
		if err != nil {
			errorEncode(w, infrastruct.ErrBadJSONURL)
			return
		}
	}

	ctx := r.Context()
	outputMultiplexer, err := h.service.SrvMultiplexer(ctx, arrURLs.URLs)
	if err != nil {
		errorEncode(w, err)
		return
	}

	responseEncoder(w, outputMultiplexer)
}
