package modules

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/heroiclabs/nakama-common/runtime"
)

type Notification struct {
	Message string `json:"Message"`
}

type InGameNotification struct {
	MatchId string `json:"MatchId"`
	Message string `json:"Message"`
}

func GameNotification(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	subject := payload

	content := map[string]interface{}{}

	code := 101

	persistent := false
	err := nk.NotificationSendAll(ctx, subject, content, code, persistent)

	return subject, err
}

func InGameNotificationMatch(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	var req map[string]interface{}

	jsonErr := json.Unmarshal([]byte(payload), &req)

	if jsonErr != nil {
		logger.Error("SendInGameNotiToOneMatch Unmarshal Error")
	}

	matchId := req["MatchId"].(string)
	subject := req["Message"].(string)

	presences := GetPresencesByMathId(matchId)

	code := 101
	content := map[string]interface{}{}

	persistent := false

	gameNotification := Notification{subject}
	jsonData, _ := json.Marshal(gameNotification)

	for _, p := range presences {
		logger.Debug(p.GetUserId())
		nk.NotificationSend(ctx, p.GetUserId(), string(jsonData), content, code, "", persistent)
	}

	return "", nil
}
