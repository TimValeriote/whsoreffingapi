package models

type Season struct {
	Id           int
	SeasonName   string
	SeasonTypeId int
	LeagueId     int
	LsSeasonId   int
	Active       bool
}

type SeasonService interface {
	GetAllActiveSeasons() ([]Season, error)
	GetSeasonById(seasonId int) (Season, error)
}
