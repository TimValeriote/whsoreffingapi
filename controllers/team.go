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

type TeamController struct {
	Log *logrus.Logger
}

type TeamResponse struct {
	Teams []TeamInfoStruct `json:"teams"`
}

type TeamInfoStruct struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	NickName string `json:"nickname"`
	LeagueId int    `json:"league_id"`
}
