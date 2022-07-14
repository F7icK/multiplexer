package config

import (
	"github.com/F7icK/multiplexer/internal/multiplexer/types"
)

func NewConfig() *types.Config {
	return &types.Config{
		LimitConnection: 100,
		LimitGoRoutines: 4,
		TimeoutOutgoing: 1,
		TimeoutIncoming: 10,
		Port:            ":8080",
	}
}
