package types

import "github.com/heroiclabs/nakama-common/runtime"

const (
	MatchStatusNotStarted = 0
	MatchStatusRunning    = 1
	MatchStatusFinished   = 3
)

type MatchState struct {
	Presences      map[string]runtime.Presence
	Label          *MatchLabel
	EmptyTicks     int
	Status         int
	Playing        bool
	JoinInProgress int
}
type MatchLabel struct {
	Open  int `json:"open"`
	Close int `json:"close"`
}
