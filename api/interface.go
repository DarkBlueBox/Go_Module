package api

import (
	"context"
	"database/sql"
	"matchmod/config"
	"matchmod/types"

	"github.com/heroiclabs/nakama-common/runtime"
)

type MatchInterface interface {
	Init(config *config.Config)
	Update(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state *types.MatchState, messages []runtime.MatchData)
}
