package types

import (
	"sync"
	"time"
)

type URLsRequest struct {
	URLs []string `json:"urls"`
}

type ResultBody struct {
	Results []DataWrite `json:"results"`
}

type DataWrite struct {
	URL    string      `json:"url"`
	Result interface{} `json:"result"`
}

type DataWriteChan struct {
	URL    string      `json:"url"`
	Result interface{} `json:"result"`
	Err    error       `json:"err"`
}

type Connection struct {
	Connected       uint32     `json:"сonnected"`
	LimitConnection uint32     `json:"limit_connection"`
	Mut             sync.Mutex `json:"mut"`
}

type Config struct {
	LimitConnection uint32        `json:"limit_connection"`
	LimitGoRoutines int           `json:"limit_go_routines"`
	TimeoutOutgoing time.Duration `json:"timeout_outgoing"`
	TimeoutIncoming time.Duration `json:"timeout_incoming"`
	Port            string        `json:"port"`
}
