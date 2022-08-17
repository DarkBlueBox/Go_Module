package config

type Lobby struct {
	Players map[string]*Player `json:"Players"`
}

func CreateLobby(config *Config) Lobby {
	return Lobby{
		Players: make(map[string]*Player),
	}
}

func (r *Lobby) AddPlayer(uid string, config *Config) bool {

	r.Players[uid] = &Player{
		UID: uid,
		Hp:  config.PlayerHealth,
	}
	return true
}

func (r *Lobby) DeletePlayer(uid string) {
	delete(r.Players, uid)
}
