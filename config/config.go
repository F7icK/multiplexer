package config

import "time"

type cfg struct {
	LimitConnection uint32
	LimitGoRoutines int
	TimeoutOutgoing time.Duration
	TimeoutIncoming time.Duration
}

func NewConfig() *cfg {
	return &cfg{
		LimitConnection: 100,
		LimitGoRoutines: 4,
		TimeoutOutgoing: 1,
		TimeoutIncoming: 10,
	}
}
