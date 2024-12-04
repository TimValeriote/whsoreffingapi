package models

type OfficialGame struct {
	Id         int
	GameId     int
	OfficialId int
}

type OfficialGameService interface {
	GetOfficialGamesByOfficialId(officialId int) ([]OfficialGame, error)
}
