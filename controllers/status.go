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

type StatusController struct {
	Log *logrus.Logger
}

type StatusResponse struct {
	Statuses []StatusInfoStruct `json:"statuses"`
}

type StatusInfoStruct struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
}
