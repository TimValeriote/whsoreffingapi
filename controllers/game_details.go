package controllers

import (
	"encoding/json"
	"net/http"
	"runtime/debug"
	"strconv"

	"github.com/sirupsen/logrus"
	"whosreffing-api/api/constants"
	"whosreffing-api/core/models"
	"whosreffing-api/utils"
)

type GameDetailsController struct {
	Log *logrus.Logger
}

type GameDetailsResponse struct {
	GameDetails []GameDetailsInfoStruct `json:"game_details"`
}

type GameDetailResponse struct {
	GameDetail GameDetailsInfoStruct `json:"game_details"`
}

type GameDetailsInfoStruct struct {
	Id                 int    `json:"id"`
	GameId             int    `json:"game_id"`
	GameLength         string `json:"game_length"`
	HomeGoals          int    `json:"home_goals"`
	VisitingGoals      int    `json:"visiting_goals"`
	OverTime           int    `json:"overtime"`
	Shootout           int    `json:"shootout"`
	HomePowerPlays     int    `json:"home_power_plays"`
	VisitorPowerPlays  int    `json:"visitor_power_plays"`
	HomePims           int    `json:"home_pims"`
	VisitorPims        int    `json:"visitor_pims"`
	HomeFaceOffs       int    `json:"home_faceoffs"`
	VisitorFaceOffs    int    `json:"visitor_faceoffs"`
	HomeFaceOffsWon    int    `json:"home_faceoffs_won"`
	VisitorFaceOffsWon int    `json:"visitor_faceoffs_won"`
}

func (controller GameDetailsController) GetOfficialGamesByOfficialId(writer http.ResponseWriter, request *http.Request) {
	context, err := utils.NewServiceFromContext(request, constants.CONTEXT_PARAMS, constants.CONTEXT_LOGGER, constants.CONTEXT_CORE)
	if err != nil {
		context.Log.WithFields(logrus.Fields{
			"event":      "whosreffing-api::GameDetailsController::GetGameDetailsByGameId - Failed to get value from context",
			"stackTrace": string(debug.Stack()),
		}).Error(err)
		return
	}

	officialId, err := strconv.Atoi(context.Params.ByName("officialId"))
	if err != nil {
		http.Error(writer, `{"error": "No valid id was provided"}`, http.StatusBadRequest)
		return
	}

	games, err := context.Core.GameDetailsService.GetOfficialGameDetailsByOfficialID(officialId)
	if err != nil {
		http.Error(writer, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}

	response := ConstructGameDetailsResponse(games)

	context.Core.Commit()

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(response)
}

func (controller GameDetailsController) GetGameDetailsByGameId(writer http.ResponseWriter, request *http.Request) {
	context, err := utils.NewServiceFromContext(request, constants.CONTEXT_PARAMS, constants.CONTEXT_LOGGER, constants.CONTEXT_CORE)
	if err != nil {
		context.Log.WithFields(logrus.Fields{
			"event":      "whosreffing-api::GameDetailsController::GetGameDetailsByGameId - Failed to get value from context",
			"stackTrace": string(debug.Stack()),
		}).Error(err)
		return
	}

	gameId, err := strconv.Atoi(context.Params.ByName("gameId"))
	if err != nil {
		http.Error(writer, `{"error": "No valid id was provided"}`, http.StatusBadRequest)
		return
	}

	game, err := context.Core.GameDetailsService.GetGameDetailsByGameId(gameId)
	if err != nil {
		http.Error(writer, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}

	response := ConstructGameDetailResponse(game)

	context.Core.Commit()

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(response)
}

func ConstructGameDetailResponse(arena models.GameDetails) GameDetailResponse {
	var response GameDetailResponse
	gameDetailsResponse := ConstructGameDetailsStructResponse(arena)
	response.GameDetail = gameDetailsResponse
	return response
}

func ConstructGameDetailsResponse(games []models.GameDetails) GameDetailsResponse {
	var response GameDetailsResponse
	gamesArray := make([]GameDetailsInfoStruct, 0)
	for _, game := range games {
		gameResponse := ConstructGameDetailsStructResponse(game)
		gamesArray = append(gamesArray, gameResponse)
	}
	response.GameDetails = gamesArray
	return response
}

func ConstructGameDetailsStructResponse(gameDetails models.GameDetails) GameDetailsInfoStruct {
	return GameDetailsInfoStruct{
		Id:                 gameDetails.Id,
		GameId:             gameDetails.GameId,
		GameLength:         gameDetails.GameLength,
		HomeGoals:          gameDetails.HomeGoals,
		VisitingGoals:      gameDetails.VisitingGoals,
		OverTime:           gameDetails.OverTime,
		Shootout:           gameDetails.Shootout,
		HomePowerPlays:     gameDetails.HomePowerPlays,
		VisitorPowerPlays:  gameDetails.VisitorPowerPlays,
		HomePims:           gameDetails.HomePims,
		VisitorPims:        gameDetails.VisitorPims,
		HomeFaceOffs:       gameDetails.HomeFaceOffs,
		VisitorFaceOffs:    gameDetails.VisitorFaceOffs,
		HomeFaceOffsWon:    gameDetails.HomeFaceOffsWon,
		VisitorFaceOffsWon: gameDetails.VisitorFaceOffsWon,
	}
}
