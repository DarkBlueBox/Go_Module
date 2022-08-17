package main

import (
	"context"
	"database/sql"
	"matchmod/modules"

	"github.com/heroiclabs/nakama-common/runtime"
)

func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	logger.Info("Hello World!")

	if err := initializer.RegisterMatchmakerMatched(modules.MakeMatch); err != nil {
		logger.Error("Unable to register: %v", err)
		return err
	}

	if err := initializer.RegisterRpc("create_match_rpc", modules.CreateMatchRPC); err != nil {
		logger.Error("Unable to register: %v", err)
		return err
	}

	if err := initializer.RegisterRpc("list_matches_rpc", modules.ListMatchesRPC); err != nil {
		logger.Error("Unable to register: %v", err)
		return err
	}

	if err := initializer.RegisterRpc("game_notification", modules.GameNotification); err != nil {
		logger.Error("Unable to register: %v", err)
		return err
	}

	if err := initializer.RegisterRpc("in_game_notificationM_match", modules.InGameNotificationMatch); err != nil {
		logger.Error("Unable to register: %v", err)
		return err
	}
	if err := initializer.RegisterMatch("some_match", modules.RegisterMatch); err != nil {
		logger.Error("Unable to register: %v", err)
		return err
	}

	return nil
}
