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

type OfficialGameController struct {
	Log *logrus.Logger
}

type OfficialGamesResponse struct {
	OfficialGames []OfficialGameInfoStruct `json:"official_games"`
}

type OfficialGameInfoStruct struct {
	Id         int `json:"id"`
	GameId     int `json:"game_id"`
	OfficialId int `json:"official_id"`
}

func (controller OfficialGameController) GetOfficialGamesByOfficialId(writer http.ResponseWriter, request *http.Request) {
	context, err := utils.NewServiceFromContext(request, constants.CONTEXT_PARAMS, constants.CONTEXT_LOGGER, constants.CONTEXT_CORE)
	if err != nil {
		context.Log.WithFields(logrus.Fields{
			"event":      "whosreffing-api::OfficialGameController::GetOfficialGamesByOfficialId - Failed to get value from context",
			"stackTrace": string(debug.Stack()),
		}).Error(err)
		return
	}

	officialId, err := strconv.Atoi(context.Params.ByName("officialId"))
	if err != nil {
		http.Error(writer, `{"error": "No valid id was provided"}`, http.StatusBadRequest)
		return
	}

	officialGames, err := context.Core.OfficialGameService.GetOfficialGamesByOfficialId(officialId)
	if err != nil {
		http.Error(writer, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}

	response := ConstructOfficialGamesResponse(officialGames)

	context.Core.Commit()

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(response)
}

func ConstructOfficialGamesResponse(officialGames []models.OfficialGame) OfficialGamesResponse {
	var response OfficialGamesResponse
	officialGamesArray := make([]OfficialGameInfoStruct, 0)
	for _, officialGame := range officialGames {
		officialGameResponse := ConstructOfficialGameStructResponse(officialGame)
		officialGamesArray = append(officialGamesArray, officialGameResponse)
	}
	response.OfficialGames = officialGamesArray
	return response
}

func ConstructOfficialGameStructResponse(officialGame models.OfficialGame) OfficialGameInfoStruct {
	return OfficialGameInfoStruct{
		Id:         officialGame.Id,
		GameId:     officialGame.GameId,
		OfficialId: officialGame.OfficialId,
	}
}
