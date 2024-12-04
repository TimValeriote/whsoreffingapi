package controllers

import (
	"encoding/json"
	"net/http"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"whosreffing-api/api/constants"
	"whosreffing-api/core/models"
	"whosreffing-api/utils"
)

type GameController struct {
	Log *logrus.Logger
}

type GameResponse struct {
	Game GameInfoStruct `json:"game"`
}

type GameInfoStruct struct {
	Id             int    `json:"id"`
	LsGameID       int    `json:"ls_game_id"`
	LeagueId       int    `json:"league_id"`
	Date           string `json:"date"`
	Time           string `json:"time"`
	HomeTeamId     int    `json:"home_team_id"`
	VisitingTeamId int    `json:"visiting_team_id"`
	ArenaId        int    `json:"arena_id"`
	StatusId       int    `json:"status_id"`
	RefereeOneId   int    `json:"referee_one_id"`
	RefereeTwoId   int    `json:"referee_two_id"`
	LinesmanOneId  int    `json:"linesman_one_id"`
	LinesmanTwoId  int    `json:"linesman_two_id"`
	SeasonId       int    `json:"season_id"`
}

type GamesCalendarResponse struct {
	Games []GamesCalendarInfoStruct `json:"games_calendar"`
}

type GamesCalendarInfoStruct struct {
	Date      time.Time `json:"date"`
	GameCount int       `json:"game_count"`
}

func (controller GameController) GetGamesCalendar(writer http.ResponseWriter, request *http.Request) {
	context, err := utils.NewServiceFromContext(request, constants.CONTEXT_LOGGER, constants.CONTEXT_CORE)
	if err != nil {
		context.Log.WithFields(logrus.Fields{
			"event":      "whosreffing-api::GameController::GetGamesCalendar - Failed to get value from context",
			"stackTrace": string(debug.Stack()),
		}).Error(err)
		return
	}

	games, err := context.Core.GameService.GetGamesCalendar()
	if err != nil {
		http.Error(writer, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}

	response := ConstructGamesCalendarResponse(games)

	context.Core.Commit()

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(response)
}

func ConstructGamesCalendarResponse(games []models.GamesCalendar) GamesCalendarResponse {
	var response GamesCalendarResponse
	gamesCalendarResponse := make([]GamesCalendarInfoStruct, 0)
	for _, game := range games {
		gameResponse := ConstructGamesCalendarStructResponse(game)
		gamesCalendarResponse = append(gamesCalendarResponse, gameResponse)
	}

	response.Games = gamesCalendarResponse

	return response
}

func ConstructGamesCalendarStructResponse(game models.GamesCalendar) GamesCalendarInfoStruct {
	return GamesCalendarInfoStruct{
		Date:      game.Date,
		GameCount: game.GameCount,
	}
}

func (controller GameController) GetGameById(writer http.ResponseWriter, request *http.Request) {
	context, err := utils.NewServiceFromContext(request, constants.CONTEXT_PARAMS, constants.CONTEXT_LOGGER, constants.CONTEXT_CORE)
	if err != nil {
		context.Log.WithFields(logrus.Fields{
			"event":      "whosreffing-api::GameController::GetGameById - Failed to get value from context",
			"stackTrace": string(debug.Stack()),
		}).Error(err)
		return
	}

	gameId, err := strconv.Atoi(context.Params.ByName("gameId"))
	if err != nil {
		http.Error(writer, `{"error": "No valid id was provided"}`, http.StatusBadRequest)
		return
	}

	game, err := context.Core.GameService.GetGameById(gameId)
	if err != nil {
		http.Error(writer, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}

	response := ConstructGameResponse(game)

	context.Core.Commit()

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(response)
}

func ConstructGameResponse(game models.Game) GameResponse {
	var response GameResponse
	gameResponse := ConstructGameStructResponse(game)
	response.Game = gameResponse
	return response
}

func ConstructGameStructResponse(game models.Game) GameInfoStruct {
	return GameInfoStruct{
		Id:             game.Id,
		LsGameID:       game.LsGameID,
		LeagueId:       game.LeagueId,
		Date:           game.Date,
		Time:           game.Time,
		HomeTeamId:     game.HomeTeamId,
		VisitingTeamId: game.VisitingTeamId,
		ArenaId:        game.ArenaId,
		StatusId:       game.StatusId,
		RefereeOneId:   game.RefereeOneId,
		RefereeTwoId:   game.RefereeTwoId,
		LinesmanOneId:  game.LinesmanOneId,
		LinesmanTwoId:  game.LinesmanTwoId,
		SeasonId:       game.SeasonId,
	}
}
