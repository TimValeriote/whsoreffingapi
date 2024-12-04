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

type LeagueController struct {
	Log *logrus.Logger
}

type LeaguesResponse struct {
	Leagues []LeagueInfoStruct `json:"leagues"`
}

type LeagueResponse struct {
	League LeagueInfoStruct `json:"league"`
}

type LeagueInfoStruct struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	LevelId int    `json:"level_id"`
	Code    string `json:"code"`
}

type LeagueStandingsResponse struct {
	Standings []LeagueStandingsInfoStruct `json:"standings"`
}

type LeagueStandingsInfoStruct struct {
	OfficialId int    `json:"official_id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	GameCount  int    `json:"game_count"`
}

func (controller LeagueController) GetAllLeagues(writer http.ResponseWriter, request *http.Request) {
	context, err := utils.NewServiceFromContext(request, constants.CONTEXT_LOGGER, constants.CONTEXT_CORE)
	if err != nil {
		context.Log.WithFields(logrus.Fields{
			"event":      "whosreffing-api::LeagueController::GetAllLeagues - Failed to get value from context",
			"stackTrace": string(debug.Stack()),
		}).Error(err)
		return
	}

	leagues, err := context.Core.LeagueService.GetAllLeagues()
	if err != nil {
		http.Error(writer, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}

	response := ConstructLeaguesResponse(leagues)

	context.Core.Commit()

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(response)
}

func (controller LeagueController) GetLeagueById(writer http.ResponseWriter, request *http.Request) {
	context, err := utils.NewServiceFromContext(request, constants.CONTEXT_PARAMS, constants.CONTEXT_LOGGER, constants.CONTEXT_CORE)
	if err != nil {
		context.Log.WithFields(logrus.Fields{
			"event":      "whosreffing-api::LeagueController::GetLeagueById - Failed to get value from context",
			"stackTrace": string(debug.Stack()),
		}).Error(err)
		return
	}

	leagueId, err := strconv.Atoi(context.Params.ByName("leagueId"))
	if err != nil {
		http.Error(writer, `{"error": "No valid id was provided"}`, http.StatusBadRequest)
		return
	}

	league, err := context.Core.LeagueService.GetLeagueById(leagueId)
	if err != nil {
		http.Error(writer, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}

	response := ConstructLeagueResponse(league)

	context.Core.Commit()

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(response)
}

func (controller LeagueController) GetLeagueStandingsByLeagueId(writer http.ResponseWriter, request *http.Request) {
	context, err := utils.NewServiceFromContext(request, constants.CONTEXT_PARAMS, constants.CONTEXT_LOGGER, constants.CONTEXT_CORE)
	if err != nil {
		context.Log.WithFields(logrus.Fields{
			"event":      "whosreffing-api::LeagueController::GetLeagueById - Failed to get value from context",
			"stackTrace": string(debug.Stack()),
		}).Error(err)
		return
	}

	leagueId, err := strconv.Atoi(context.Params.ByName("leagueId"))
	if err != nil {
		http.Error(writer, `{"error": "No valid id was provided"}`, http.StatusBadRequest)
		return
	}

	standings, err := context.Core.LeagueService.GetLeagueStandingsByLeagueId(leagueId)
	if err != nil {
		http.Error(writer, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}

	response := ConstructLeagueStandingsResponse(standings)

	context.Core.Commit()

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(response)
}

func ConstructLeagueStandingsResponse(standings []models.LeagueStanding) LeagueStandingsResponse {
	var response LeagueStandingsResponse
	standingsArray := make([]LeagueStandingsInfoStruct, 0)
	for _, indiv := range standings {
		standingsResponse := ConstructLeagueStandingsInfoStruct(indiv)
		standingsArray = append(standingsArray, standingsResponse)
	}
	response.Standings = standingsArray
	return response
}

func ConstructLeagueStandingsInfoStruct(indiv models.LeagueStanding) LeagueStandingsInfoStruct {
	return LeagueStandingsInfoStruct{
		OfficialId: indiv.OfficialId,
		FirstName:  indiv.FirstName,
		LastName:   indiv.LastName,
		GameCount:  indiv.GameCount,
	}
}

func ConstructLeaguesResponse(leagues []models.League) LeaguesResponse {
	var response LeaguesResponse
	leagueArray := make([]LeagueInfoStruct, 0)
	for _, league := range leagues {
		leagueResponse := ConstructLeagueStructResponse(league)
		leagueArray = append(leagueArray, leagueResponse)
	}
	response.Leagues = leagueArray
	return response
}

func ConstructLeagueResponse(league models.League) LeagueResponse {
	var response LeagueResponse
	leagueResponse := ConstructLeagueStructResponse(league)
	response.League = leagueResponse
	return response
}

func ConstructLeagueStructResponse(league models.League) LeagueInfoStruct {
	return LeagueInfoStruct{
		Id:      league.Id,
		Name:    league.Name,
		LevelId: league.LevelId,
		Code:    league.Code,
	}
}
