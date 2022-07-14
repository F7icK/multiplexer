package handlers

import (
	"encoding/json"
	"errors"
	"github.com/F7icK/multiplexer/internal/multiplexer/service"
	"github.com/F7icK/multiplexer/internal/multiplexer/types"
	"github.com/F7icK/multiplexer/pkg/infrastruct"
	"log"
	"net/http"
	"net/url"
	"syscall"
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
		log.Println(err)
	}
}

func responseEncoder(w http.ResponseWriter, res interface{}) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)

	if err := json.NewEncoder(w).Encode(res); err != nil && !errors.Is(err, syscall.EPIPE) {
		log.Println(err)
	}
}

func (h *Handlers) PostURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	jsonUrls := types.UrlsRequest{}
	if err := json.NewDecoder(r.Body).Decode(&jsonUrls); err != nil {
		errorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	countUrl := len(jsonUrls.Urls)
	if countUrl == 0 || countUrl > 20 {
		errorEncode(w, infrastruct.ErrorCountUrl)
		return
	}

	for _, urlInUrls := range jsonUrls.Urls {
		_, err := url.ParseRequestURI(urlInUrls)
		if err != nil {
			errorEncode(w, infrastruct.ErrorBadJsonUrl)
			return
		}
	}

	outputData, err := h.service.GetUrls(ctx, jsonUrls.Urls)
	if err != nil {
		errorEncode(w, err)
		return
	}

	responseEncoder(w, outputData)
}
