package types

import (
	"sync"
	"time"
)

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

type DataWriteChan struct {
	Url    string
	Result interface{}
	Error  error
}

type Connection struct {
	Connected       uint32     `json:"сonnected"`
	LimitConnection uint32     `json:"limit_connection"`
	Mut             sync.Mutex `json:"mut"`
}

type Config struct {
	LimitConnection uint32
	LimitGoRoutines int
	TimeoutOutgoing time.Duration
	TimeoutIncoming time.Duration
	Port            string
}
