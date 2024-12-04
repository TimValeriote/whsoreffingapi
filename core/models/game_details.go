package models

type GameDetails struct {
	Id                 int
	GameId             int
	GameLength         string
	HomeGoals          int
	VisitingGoals      int
	OverTime           int
	Shootout           int
	HomePowerPlays     int
	VisitorPowerPlays  int
	HomePims           int
	VisitorPims        int
	HomeFaceOffs       int
	VisitorFaceOffs    int
	HomeFaceOffsWon    int
	VisitorFaceOffsWon int
}

type GameDetailsService interface {
	GetGameDetailsByGameId(gameId int) (GameDetails, error)
	GetOfficialGameDetailsByOfficialID(officialId int) ([]GameDetails, error)
}
