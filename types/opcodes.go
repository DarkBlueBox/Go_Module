package types

const (
	PlayerJoined = 1
	PlayerLeft   = 2
)

type PlayerUpdateLeft struct {
	UID string
}
