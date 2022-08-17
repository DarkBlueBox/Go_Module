package modules

import (
	"context"
	"database/sql"
	"matchmod/api"
	"matchmod/config"

	"github.com/heroiclabs/nakama-common/runtime"
)

func MakeMatch(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, entries []runtime.MatchmakerEntry) (string, error) {
	for _, e := range entries {
		logger.Info("Matched user '%s' named '%s'", e.GetPresence().GetUserId(), e.GetPresence().GetUsername())

		for k, v := range e.GetProperties() {
			logger.Info("Matched on '%s' value '%v'", k, v)
		}
	}

	matchID, err := nk.MatchCreate(ctx, "somematch", map[string]interface{}{"invited": entries})
	if err != nil {
		return "", err
	}

	return matchID, nil
}

func RegisterMatch(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule) (runtime.Match, error) {
	match := []api.MatchInterface{
		&api.Msghandler{},
	}

	matchConfig := config.Config{
		PlayerHealth: 100,
		PlayerSpeed:  1,
	}

	return &MatchHandler{api: match, config: matchConfig}, nil
}
