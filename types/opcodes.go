package types

const (
	LobbyState   = 0
	PlayerJoined = 1
	PlayerLeft   = 2
	PlayerMove   = 3
	MatchPause   = 4
	MatchUnpause = 5
	MatchEnd     = 6
)

type Message struct {
	Ok      bool
	Message string
}

type PlayerUpdateLeft struct {
	UID string
}

type PlayerUpdateMove struct {
	UID    string
	From   int
	To     int
	Result Message
}
type PlayerDataMove struct {
	Point int
}
