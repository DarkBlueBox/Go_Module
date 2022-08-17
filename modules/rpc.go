package modules

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/heroiclabs/nakama-common/runtime"
)

func CreateMatchRPC(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	logger.Info("Payload: %s", payload)

	params := make(map[string]interface{})
	if err := json.Unmarshal([]byte(payload), &params); err != nil {
		return "", err
	}

	modulename := "somematch"
	matchID, err := nk.MatchCreate(ctx, modulename, params)
	if err != nil {
		return "", err
	}
	return matchID, nil
}

func ListMatchesRPC(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	limit := 10
	authoritative := true
	label := ""
	minSize := 0
	maxSize := 2
	query := "+label.open:1 +label.close:1"
	matches, _ := nk.MatchList(ctx, limit, authoritative, label, &minSize, &maxSize, query)

	data, _ := json.Marshal(matches)

	return string(data), nil
}
