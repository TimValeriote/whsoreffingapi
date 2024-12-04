package models

type Team struct {
	Id       int
	Name     string
	Code     string
	NickName string
	LeagueId int
}

type TeamService interface {
	GetTeamById(teamId int) (Team, error)
}
