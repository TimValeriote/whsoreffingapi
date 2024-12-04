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

type GameTypeController struct {
	Log *logrus.Logger
}

type GameTypeResponse struct {
	GameTypes []GameTypeInfoStruct `json:"game_types"`
}

type GameTypeInfoStruct struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
}
