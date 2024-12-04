package controllers

import (
	//"encoding/json"
	//"net/http"
	//"runtime/debug"
	//"strconv"

	"github.com/sirupsen/logrus"
	//"phl-skate-sharpening-api/api/constants"
	//"phl-skate-sharpening-api/core/models"
	//"phl-skate-sharpening-api/utils"
)

type SeasonController struct {
	Log *logrus.Logger
}

type SeasonResponse struct {
	Seasons []SeasonInfoStruct `json:"seasons"`
}

type SeasonInfoStruct struct {
	Id           int    `json:"id"`
	SeasonName   string `json:"season_name"`
	SeasonTypeId int    `json:"season_type_id"`
	LeagueId     int    `json:"league_id"`
	LsSeasonId   int    `json:"ls_season_ud"`
	Active       bool   `json:"active"`
}
