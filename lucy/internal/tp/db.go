package tp

import "github.com/mini-projects/keylogger-server/internal/database"

type ApiConf struct {
	DB *database.Queries
}

type CounterParameters struct {
	Counter int64 `json:"counter"`
}

type MouseStats struct {
	MouseDistance int64 `json:"MouseDistance"`
	LeftClick int64 `json:"leftClick"`
	RightClick int64 `json:"rightClick"`
}
