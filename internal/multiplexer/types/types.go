package types

import "sync"

type UrlsRequest struct {
	Urls []string `json:"urls"`
}

type ResultBody struct {
	Results []DataWrite `json:"results"`
}

type DataWrite struct {
	Url    string      `json:"url"`
	Result interface{} `json:"result"`
}

type Connection struct {
	Connected       uint32     `json:"сonnected"`
	LimitConnection uint32     `json:"limit_connection"`
	Mut             sync.Mutex `json:"mut"`
}
