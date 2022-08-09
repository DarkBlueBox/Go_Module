package modules

import (
	"context"
	"database/sql"
	"encoding/json"
	"matchmod/types"

	"github.com/heroiclabs/nakama-common/runtime"
	"google.golang.org/protobuf/encoding/protojson"
)

const newMatchOpCode = 999

type MatchHandler struct {
	Marshaler   *protojson.MarshalOptions
	Unmarshaler *protojson.UnmarshalOptions
}

var _matchMap = map[string]types.MatchState{}

func GetPresencesByMathId(matchId string) map[string]runtime.Presence {
	matchState := _matchMap[matchId]
	return matchState.Presences
}

func (m *MatchHandler) MatchInit(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, params map[string]interface{}) (interface{}, int, string) {
	close, ok := params["close"].(bool)
	if !ok {
		logger.Error("invalid match init parameter \"close\"")
		return nil, 0, ""
	}

	tickRate := 20
	label := &types.MatchLabel{
		Open: 1,
	}
	if close {
		label.Close = 1
	}

	labelJSON, err := json.Marshal(label)
	if err != nil {
		logger.WithField("error", err).Error("match init failed")
		labelJSON = []byte("{}")
	}
	state := &types.MatchState{
		EmptyTicks: 0,
		Label:      label,
		Presences:  map[string]runtime.Presence{},
		Status:     types.MatchStatusNotStarted,
	}

	return state, tickRate, string(labelJSON)
}
func (m *MatchHandler) MatchJoinAttempt(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presence runtime.Presence, metadata map[string]string) (interface{}, bool, string) {
	result := true
	matchState := state.(*types.MatchState)

	if presence, ok := matchState.Presences[presence.GetUserId()]; ok {
		if presence == nil {
			matchState.JoinInProgress++
			return matchState, true, ""
		} else {
			return matchState, false, "already joined"
		}
	}

	if len(matchState.Presences)+matchState.JoinInProgress >= 2 {
		return matchState, false, "match full"
	}
	matchState.JoinInProgress++

	return matchState, result, ""
}

func (m *MatchHandler) MatchJoin(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presences []runtime.Presence) interface{} {
	matchState, ok := state.(*types.MatchState)
	if !ok {
		logger.Error("state not a valid lobby state object")
	}

	for _, presence := range matchState.Presences {
		matchState.EmptyTicks = 0
		matchState.Presences[presence.GetUserId()] = presence
		matchState.Presences[presence.GetSessionId()] = presence
		matchState.JoinInProgress--

		if len(matchState.Presences) >= 2 && matchState.Label.Open != 0 {
			matchState.Label.Open = 0
			if labelJSON, err := json.Marshal(matchState.Label); err != nil {
				logger.Error("error encoding label: %v", err)
			} else {
				if err := dispatcher.MatchLabelUpdate(string(labelJSON)); err != nil {
					logger.Error("error updating label: %v", err)
				}
			}
		}
		data, err := json.Marshal(matchState)
		if err != nil {
			logger.Error("Could not json.Marshal the state.")
		}
		dispatcher.BroadcastMessage(types.PlayerJoined, data, nil, nil, true)
	}

	return matchState
}

func (m *MatchHandler) MatchLeave(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presences []runtime.Presence) interface{} {
	matchState, ok := state.(*types.MatchState)
	if !ok {
		logger.Error("state not a valid lobby state object")
	}

	for _, p := range presences {
		logger.Info("Player %v left.", p.GetUserId())
		delete(matchState.Presences, p.GetUserId())
		data, err := json.Marshal(types.PlayerUpdateLeft{UID: p.GetUserId()})
		if err != nil {
			logger.Error("Could not json.Marshal.")
		}
		dispatcher.BroadcastMessage(types.PlayerLeft, data, nil, nil, true)
	}

	return matchState
}

func (m *MatchHandler) MatchLoop(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, messages []runtime.MatchData) interface{} {

	matchState, ok := state.(*types.MatchState)
	if !ok {
		logger.Error("state not a valid lobby state object")
	}

	if len(matchState.Presences) == 0 {
		matchState.EmptyTicks++
	}

	if matchState.EmptyTicks > 100 {
		return nil
	}

	return matchState
}

func (m *MatchHandler) MatchSignal(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, data string) (interface{}, string) {
	return state, "signal received: " + data
}

func (m *MatchHandler) MatchTerminate(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, graceSeconds int) interface{} {
	logger.Debug("match will terminate in %d seconds", graceSeconds)

	var matchId string

	// Find an existing match for the remaining connected presences to join
	limit := 1
	authoritative := true
	label := ""
	minSize := 2
	maxSize := 4
	query := "*"
	availableMatches, err := nk.MatchList(ctx, limit, authoritative, label, &minSize, &maxSize, query)
	if err != nil {
		logger.Error("error listing matches", err)
		return nil
	}

	if len(availableMatches) > 0 {
		matchId = availableMatches[0].MatchId
	} else {
		// No available matches, create a new match instead
		matchId, err = nk.MatchCreate(ctx, "match", nil)
		if err != nil {
			logger.Error("error creating match", err)
			return nil
		}
	}

	// Broadcast the new match id to all remaining connected presences
	data := map[string]string{
		matchId: matchId,
	}

	dataJson, err := json.Marshal(data)
	if err != nil {
		logger.Error("error marshaling new match message")
		return nil
	}

	dispatcher.BroadcastMessage(newMatchOpCode, dataJson, nil, nil, true)

	return state
}
