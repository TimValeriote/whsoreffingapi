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

type ArenaController struct {
	Log *logrus.Logger
}

type ArenaResponse struct {
	Arena ArenaInfoStruct `json:"arena"`
}

type ArenaInfoStruct struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (controller ArenaController) GetArenaById(writer http.ResponseWriter, request *http.Request) {
	context, err := utils.NewServiceFromContext(request, constants.CONTEXT_PARAMS, constants.CONTEXT_LOGGER, constants.CONTEXT_CORE)
	if err != nil {
		context.Log.WithFields(logrus.Fields{
			"event":      "whosreffing-api::ArenaController::GetArenaById - Failed to get value from context",
			"stackTrace": string(debug.Stack()),
		}).Error(err)
		return
	}

	arenaId, err := strconv.Atoi(context.Params.ByName("arenaId"))
	if err != nil {
		http.Error(writer, `{"error": "No valid id was provided"}`, http.StatusBadRequest)
		return
	}

	arenas, err := context.Core.ArenaService.GetArenaById(arenaId)
	if err != nil {
		http.Error(writer, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}

	response := ConstructArenaResponse(arenas)

	context.Core.Commit()

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(response)
}

func ConstructArenaResponse(arena models.Arena) ArenaResponse {
	var response ArenaResponse
	arenaResponse := ConstructArenaStructResponse(arena)
	response.Arena = arenaResponse
	return response
}

func ConstructArenaStructResponse(arena models.Arena) ArenaInfoStruct {
	return ArenaInfoStruct{
		Id:   arena.Id,
		Name: arena.Name,
	}
}
