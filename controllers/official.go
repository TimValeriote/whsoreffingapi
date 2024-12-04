package controllers

import (
	"encoding/json"
	"math"
	"net/http"
	"runtime/debug"
	"strconv"

	"github.com/sirupsen/logrus"
	"whosreffing-api/api/constants"
	"whosreffing-api/core/models"
	"whosreffing-api/utils"
)

type OfficialController struct {
	Log *logrus.Logger
}

type OfficialResponse struct {
	Official OfficialInfoStruct `json:"official"`
}

type OfficialsResponse struct {
	Officials []OfficialSearchInfoStruct `json:"officials"`
}

type OfficialInfoStruct struct {
	Id           int    `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	JerseyNumber int    `json:"jersey_number"`
}

type OfficialSearchInfoStruct struct {
	Id           int    `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	JerseyNumber int    `json:"jersey_number"`
	Leagues      string `json:"leagues"`
}

type OfficialDetailsResponse struct {
	Official OfficialDetails `json:"official_details"`
}

type OfficialDetails struct {
	Id            int                     `json:"id"`
	FirstName     string                  `json:"first_name"`
	LastName      string                  `json:"last_name"`
	JerseyNumber  int                     `json:"jersey_number"`
	Games         []OfficialGamesByLeague `json:"games"`
	OfficialStats OfficialStats           `json:"official_stats"`
}

type OfficialGamesByLeague struct {
	League   string         `json:"league_name"`
	LeagueId string         `json:"league_id"`
	Games    []OfficialGame `json:"league_games"`
}

type OfficialGame struct {
	GameId      string `json:"game_id"`
	LsGameId    string `json:"ls_game_id"`
	LeagueName  string `json:"league_name"`
	LeagueId    string `json:"league_id"`
	HomeTeam    string `json:"home_team"`
	VisitorTeam string `json:"visitor_team"`
	Arena       string `json:"arena"`
	GameTime    string `json:"game_time"`
	GameDate    string `json:"game_date"`
	RefereeOne  string `json:"referee_one"`
	RefereeTwo  string `json:"referee_two"`
	LinesmanOne string `json:"linesman_one"`
	LinesmanTwo string `json:"linesman_two"`
}

type OfficialStats struct {
	Leagues                         []string `json:"leagues"`
	TotalGames                      int      `json:"total_games"`
	OvertimeAverage                 float32  `json:"overtime_average"`
	ShootoutAverage                 float32  `json:"shootout_average"`
	AverageGameTime                 string   `json:"average_game_time"`
	AverageGoalsPerGame             float32  `json:"average_goals_per_game"`
	AverageHomeGoalsPerGame         float32  `json:"average_home_goals_per_game"`
	AverageVisitorGoalsPerGame      float32  `json:"average_visitor_goals_per_game"`
	AverageHomePowerPlaysPerGame    float32  `json:"average_home_penalties_per_game"`
	AverageVisitorPowerPlaysPerGame float32  `json:"average_visitor_penalties_per_game"`
	AverageHomePimsPerGame          float32  `json:"average_home_pims_per_game"`
	AverageVisitorPimsPerGame       float32  `json:"average_visitor_pims_per_game"`
}

func (controller OfficialController) GetOfficialsBySearchTerm(writer http.ResponseWriter, request *http.Request) {
	context, err := utils.NewServiceFromContext(request, constants.CONTEXT_PARAMS, constants.CONTEXT_LOGGER, constants.CONTEXT_CORE)
	if err != nil {
		context.Log.WithFields(logrus.Fields{
			"event":      "whosreffing-api::OfficialController::GetOfficialsBySearchTerm - Failed to get value from context",
			"stackTrace": string(debug.Stack()),
		}).Error(err)
		return
	}

	// Retrieve the "searchTerm" from the query parameters
	searchTerm := request.URL.Query().Get("searchTerm")

	officials, err := context.Core.OfficialService.GetOfficialsBySearchTerm(searchTerm)
	if err != nil {
		http.Error(writer, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}

	response := ConstructOfficialsResponse(officials)

	context.Core.Commit()

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(response)
}

func (controller OfficialController) GetOfficialDetailsById(writer http.ResponseWriter, request *http.Request) {
	context, err := utils.NewServiceFromContext(request, constants.CONTEXT_PARAMS, constants.CONTEXT_LOGGER, constants.CONTEXT_CORE)
	if err != nil {
		context.Log.WithFields(logrus.Fields{
			"event":      "whosreffing-api::OfficialController::GetOfficialDetailsById - Failed to get value from context",
			"stackTrace": string(debug.Stack()),
		}).Error(err)
		return
	}

	officialId, err := strconv.Atoi(context.Params.ByName("officialId"))
	if err != nil {
		http.Error(writer, `{"error": "No valid id was provided"}`, http.StatusBadRequest)
		return
	}

	officials, err := context.Core.OfficialService.GetOfficialInfoByOfficialId(officialId)
	if err != nil {
		http.Error(writer, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}

	response := ConstructOfficialDetailsResponse(officials)

	context.Core.Commit()

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(response)
}

func ConstructOfficialDetailsResponse(officialDetails models.OfficialInfo) OfficialDetailsResponse {
	var response OfficialDetailsResponse

	var gamesResponse []OfficialGamesByLeague

	for _, leagueDetails := range officialDetails.Games {
		var gamesByLeague OfficialGamesByLeague
		var allGames []OfficialGame

		gamesByLeague.League = leagueDetails.League
		gamesByLeague.LeagueId = strconv.Itoa(leagueDetails.LeagueId)

		for _, game := range leagueDetails.Games {
			allGames = append(allGames, ConstructOfficialGame(game))
		}

		gamesByLeague.Games = allGames

		gamesResponse = append(gamesResponse, gamesByLeague)
	}

	stats := ConstructOfficialStats(officialDetails.OfficialStats)
	officialResponse := ConstructOfficialDetailsStructResponse(officialDetails, gamesResponse, stats)
	response.Official = officialResponse
	return response
}

func ConstructOfficialGame(game models.OfficialGameDetails) OfficialGame {
	return OfficialGame{
		GameId:      game.GameId,
		LsGameId:    game.LsGameId,
		LeagueName:  game.LeagueName,
		LeagueId:    strconv.Itoa(game.LeagueId),
		HomeTeam:    game.HomeTeam,
		VisitorTeam: game.VisitorTeam,
		Arena:       game.Arena,
		GameTime:    game.GameTime,
		GameDate:    game.GameDate,
		RefereeOne:  game.RefereeOne,
		RefereeTwo:  game.RefereeTwo,
		LinesmanOne: game.LinesmanOne,
		LinesmanTwo: game.LinesmanTwo,
	}
}

func ConstructOfficialStats(stats models.OfficialStats) OfficialStats {
	return OfficialStats{
		Leagues:                         stats.Leagues,
		TotalGames:                      stats.TotalGames,
		OvertimeAverage:                 roundToTwoDecimals(stats.OvertimeAverage),
		ShootoutAverage:                 roundToTwoDecimals(stats.ShootoutAverage),
		AverageGameTime:                 stats.AverageGameTime,
		AverageGoalsPerGame:             roundToTwoDecimals(stats.AverageGoalsPerGame),
		AverageHomeGoalsPerGame:         roundToTwoDecimals(stats.AverageHomeGoalsPerGame),
		AverageVisitorGoalsPerGame:      roundToTwoDecimals(stats.AverageVisitorGoalsPerGame),
		AverageHomePowerPlaysPerGame:    roundToTwoDecimals(stats.AverageHomePowerPlaysPerGame),
		AverageVisitorPowerPlaysPerGame: roundToTwoDecimals(stats.AverageVisitorPowerPlaysPerGame),
		AverageHomePimsPerGame:          roundToTwoDecimals(stats.AverageHomePimsPerGame),
		AverageVisitorPimsPerGame:       roundToTwoDecimals(stats.AverageVisitorPimsPerGame),
	}
}

func ConstructOfficialDetailsStructResponse(official models.OfficialInfo, leagueGames []OfficialGamesByLeague, stats OfficialStats) OfficialDetails {
	return OfficialDetails{
		Id:            official.Id,
		FirstName:     official.FirstName,
		LastName:      official.LastName,
		JerseyNumber:  official.JerseyNumber,
		Games:         leagueGames,
		OfficialStats: stats,
	}
}

func (controller OfficialController) GetOfficialById(writer http.ResponseWriter, request *http.Request) {
	context, err := utils.NewServiceFromContext(request, constants.CONTEXT_PARAMS, constants.CONTEXT_LOGGER, constants.CONTEXT_CORE)
	if err != nil {
		context.Log.WithFields(logrus.Fields{
			"event":      "whosreffing-api::OfficialController::GetOfficialById - Failed to get value from context",
			"stackTrace": string(debug.Stack()),
		}).Error(err)
		return
	}

	officialId, err := strconv.Atoi(context.Params.ByName("officialId"))
	if err != nil {
		http.Error(writer, `{"error": "No valid id was provided"}`, http.StatusBadRequest)
		return
	}

	officials, err := context.Core.OfficialService.GetOfficialById(officialId)
	if err != nil {
		http.Error(writer, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}

	response := ConstructOfficialResponse(officials)

	context.Core.Commit()

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(response)
}

func ConstructOfficialsResponse(officials []models.OfficialSearch) OfficialsResponse {
	var response OfficialsResponse

	var allOfficials []OfficialSearchInfoStruct

	for _, official := range officials {
		official := ConstructOfficialSearchStructResponse(official)

		allOfficials = append(allOfficials, official)
	}

	response.Officials = allOfficials
	return response
}

func ConstructOfficialSearchStructResponse(official models.OfficialSearch) OfficialSearchInfoStruct {
	return OfficialSearchInfoStruct{
		Id:           official.Id,
		FirstName:    official.FirstName,
		LastName:     official.LastName,
		JerseyNumber: official.JerseyNumber,
		Leagues:      official.Leagues,
	}
}

func ConstructOfficialResponse(official models.Official) OfficialResponse {
	var response OfficialResponse
	officialResponse := ConstructOfficialStructResponse(official)
	response.Official = officialResponse
	return response
}

func ConstructOfficialStructResponse(official models.Official) OfficialInfoStruct {
	return OfficialInfoStruct{
		Id:           official.Id,
		FirstName:    official.FirstName,
		LastName:     official.LastName,
		JerseyNumber: official.JerseyNumber,
	}
}

func roundToTwoDecimals(value float32) float32 {
	return float32(math.Round(float64(value)*100) / 100)
}
