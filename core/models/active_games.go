package models

type ActiveGame struct {
	Id       int
	GameId   int
	LsGameId int
}

type ActiveGameDetails struct {
	GameId         string
	LsGameId       string
	LeagueName     string
	HomeTeam       string
	VisitorTeam    string
	Arena          string
	GameTime       string
	RefereeOne     string
	RefereeTwo     string
	LinesmanOne    string
	LinesmanTwo    string
	HasNoOfficials bool
}

type ActiveGamesService interface {
	GetActiveGames() ([]ActiveGame, error)
	GetActiveGameDetails() ([]ActiveGameDetails, error)
}
