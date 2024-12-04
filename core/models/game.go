package models

import (
	"time"
)

type Game struct {
	Id             int
	LsGameID       int
	LeagueId       int
	Date           string
	Time           string
	HomeTeamId     int
	VisitingTeamId int
	ArenaId        int
	StatusId       int
	RefereeOneId   int
	RefereeTwoId   int
	LinesmanOneId  int
	LinesmanTwoId  int
	SeasonId       int
}

type GamesCalendar struct {
	Date      time.Time
	GameCount int
}

type GameService interface {
	GetGameById(gameId int) (Game, error)
	GetGamesCalendar() ([]GamesCalendar, error)
}
