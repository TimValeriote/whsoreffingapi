package models

type Key struct {
	Id          int
	Description string
	LeagueId    int
}

type KeysService interface {
	GetAllKeys() ([]Key, error)
	GetKeyByLeagueId(leagueId int) (Key, error)
}
