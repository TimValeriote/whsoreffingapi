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

type KeysController struct {
	Log *logrus.Logger
}

type KeysResponse struct {
	Keys []KeysInfoStruct `json:"keys"`
}

type KeysInfoStruct struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	LeagueId    int    `json:"league_id"`
}
