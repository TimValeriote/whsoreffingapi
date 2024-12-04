package controllers

import (
	"encoding/json"
	"net/http"
	"runtime/debug"
	//"strconv"

	"github.com/sirupsen/logrus"
	"whosreffing-api/api/constants"
	"whosreffing-api/core/models"
	"whosreffing-api/utils"
)

type ActiveGamesController struct {
	Log *logrus.Logger
}

type ActiveGamesResponse struct {
	ActiveGames []ActiveGameInfoStruct `json:"active_games"`
}

type ActiveGameInfoStruct struct {
	Id       int `json:"id"`
	GameId   int `json:"game_id"`
	LsGameId int `json:"ls_game_id"`
}

type ActiveGamesDetailsResponse struct {
	Leagues map[string][]ActiveGameDetailsInfoStruct `json:"active_games"`
}

type ActiveGameDetailsInfoStruct struct {
	GameId         string `json:"game_id"`
	LsGameId       string `json:"ls_game_id"`
	HomeTeam       string `json:"home_team"`
	VisitorTeam    string `json:"visitor_team"`
	Arena          string `json:"arena"`
	GameTime       string `json:"game_time"`
	RefereeOne     string `json:"referee_one"`
	RefereeTwo     string `json:"referee_two"`
	LinesmanOne    string `json:"linesman_one"`
	LinesmanTwo    string `json:"linesman_two"`
	HasNoOfficials bool   `json:"hasNoOfficials"`
}

func (controller ActiveGamesController) GetActiveGames(writer http.ResponseWriter, request *http.Request) {
	context, err := utils.NewServiceFromContext(request, constants.CONTEXT_LOGGER, constants.CONTEXT_CORE)
	if err != nil {
		context.Log.WithFields(logrus.Fields{
			"event":      "whosreffing-api::ActiveGamesController::GetActiveGames - Failed to get value from context",
			"stackTrace": string(debug.Stack()),
		}).Error(err)
		return
	}

	activeGames, err := context.Core.ActiveGamesService.GetActiveGames()
	if err != nil {
		http.Error(writer, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}

	response := ConstructActiveGamesResponse(activeGames)

	context.Core.Commit()

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(response)
}

func ConstructActiveGamesResponse(activeGames []models.ActiveGame) ActiveGamesResponse {
	var response ActiveGamesResponse
	activeGamesArray := make([]ActiveGameInfoStruct, 0)
	for _, activeGame := range activeGames {
		activeGameResponse := ConstructActiveGameStructResponse(activeGame)
		activeGamesArray = append(activeGamesArray, activeGameResponse)
	}

	response.ActiveGames = activeGamesArray
	return response
}

func ConstructActiveGameStructResponse(activeGame models.ActiveGame) ActiveGameInfoStruct {
	return ActiveGameInfoStruct{
		Id:       activeGame.Id,
		GameId:   activeGame.GameId,
		LsGameId: activeGame.LsGameId,
	}
}

func (controller ActiveGamesController) GetActiveGamesDetails(writer http.ResponseWriter, request *http.Request) {
	context, err := utils.NewServiceFromContext(request, constants.CONTEXT_LOGGER, constants.CONTEXT_CORE)
	if err != nil {
		context.Log.WithFields(logrus.Fields{
			"event":      "whosreffing-api::ActiveGamesController::GetActiveGamesDetails - Failed to get value from context",
			"stackTrace": string(debug.Stack()),
		}).Error(err)
		return
	}

	activeGames, err := context.Core.ActiveGamesService.GetActiveGameDetails()
	if err != nil {
		http.Error(writer, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}

	response := ConstructActiveGamesDetailsResponse(activeGames)

	context.Core.Commit()

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(response)
}

func ConstructActiveGamesDetailsResponse(activeGames []models.ActiveGameDetails) ActiveGamesDetailsResponse {
	response := ActiveGamesDetailsResponse{
		Leagues: make(map[string][]ActiveGameDetailsInfoStruct),
	}

	for _, activeGame := range activeGames {
		activeGameResponse := ConstructActiveGameDetailsStructResponse(activeGame)

		response.Leagues[activeGame.LeagueName] = append(response.Leagues[activeGame.LeagueName], activeGameResponse)
	}

	return response
}

func ConstructActiveGameDetailsStructResponse(activeGame models.ActiveGameDetails) ActiveGameDetailsInfoStruct {
	return ActiveGameDetailsInfoStruct{
		GameId:         activeGame.GameId,
		LsGameId:       activeGame.LsGameId,
		HomeTeam:       activeGame.HomeTeam,
		VisitorTeam:    activeGame.VisitorTeam,
		Arena:          activeGame.Arena,
		GameTime:       activeGame.GameTime,
		RefereeOne:     activeGame.RefereeOne,
		RefereeTwo:     activeGame.RefereeTwo,
		LinesmanOne:    activeGame.LinesmanOne,
		LinesmanTwo:    activeGame.LinesmanTwo,
		HasNoOfficials: activeGame.HasNoOfficials,
	}
}
