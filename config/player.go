package config

type Player struct {
	UID      string `json:"UID"`
	Hp       int    `json:"HP"`
	Location int    `json:"Location"`
}
