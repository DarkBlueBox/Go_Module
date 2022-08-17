package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"matchmod/config"
	"matchmod/types"

	"github.com/heroiclabs/nakama-common/runtime"
)

type Msghandler struct {
	config *config.Config
}

func (s *Msghandler) Init(config *config.Config) {
	s.config = config
}

func (s *Msghandler) Update(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state *types.MatchState, messages []runtime.MatchData) {
	for _, message := range messages {
		op := message.GetOpCode()
		switch op {
		case types.PlayerMove:
			s.handlePlayerMove(logger, dispatcher, state, message)
			break

		}
	}
}

func (s *Msghandler) handlePlayerMove(logger runtime.Logger, dispatcher runtime.MatchDispatcher, state *types.MatchState, message runtime.MatchData) {
	uid := message.GetUserId()
	place := types.PlayerDataMove{}

	move := types.PlayerUpdateMove{
		UID:    uid,
		From:   state.Lobby.Players[uid].Location,
		To:     place.Point,
		Result: types.Message{Ok: true, Message: "OK"},
	}
	data, err := json.Marshal(move)
	if err != nil {
		logger.Error("error encoding: %v", err)
	}
	logger.Info("Player %s moved from %d to %d.", message.GetUsername(), move.From, move.To)
	dispatcher.BroadcastMessage(types.PlayerMove, data, nil, nil, true)

}
