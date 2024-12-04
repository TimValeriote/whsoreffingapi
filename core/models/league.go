package models

type League struct {
	Id      int
	Name    string
	LevelId int
	Code    string
}

type LeagueStanding struct {
	OfficialId int
	FirstName  string
	LastName   string
	GameCount  int
}

type LeagueService interface {
	GetAllLeagues() ([]League, error)
	GetLeagueById(leagueID int) (League, error)
	GetLeagueStandingsByLeagueId(leagueId int) ([]LeagueStanding, error)
}
