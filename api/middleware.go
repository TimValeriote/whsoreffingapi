package main

import (
	"database/sql"
	"net/http"
	"runtime/debug"

	"github.com/Tomasen/realip"
	"github.com/gorilla/context"
	"github.com/sirupsen/logrus"
	"whosreffing-api/api/constants"
	"whosreffing-api/core"
	"whosreffing-api/core/models"
	"whosreffing-api/utils"
)

type Middleware struct {
	Database *sql.DB
	Log      *logrus.Logger
}

func (middleware *Middleware) OptionsMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE")
		writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		writer.Header().Set("Content-Type", "application/json")

		request.RemoteAddr = realip.FromRequest(request)

		handler.ServeHTTP(writer, request)
	})
}

func (middleware *Middleware) CORSMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, HEAD, OPTIONS")
		writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		writer.Header().Set("Content-Type", "application/json")

		request.RemoteAddr = realip.FromRequest(request)

		handler.ServeHTTP(writer, request)
	})
}

func (middleware *Middleware) CoreMasterCoreMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		context.Set(request, constants.CONTEXT_USEDATABASE, constants.CONTEXT_DATABASE)
		handler.ServeHTTP(writer, request)
	})
}

func (middleware *Middleware) CoreApplicationServiceMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var db *sql.DB

		service, err := utils.NewServiceFromContext(request, constants.CONTEXT_USEDATABASE)
		if err != nil {
			service.Log.WithFields(logrus.Fields{
				"event":      "phlapi::CoreApplicationServiceMiddleware - Failed to get value from context",
				"stackTrace": string(debug.Stack()),
			}).Error(err)

			http.Error(writer, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
			return
		}

		db = middleware.Database

		co, err := core.CreateCore(db)
		if err != nil {
			service.Log.WithFields(logrus.Fields{
				"event":      "phlapi::CoreApplicationServiceMiddleware - Failed to create Core instance",
				"stackTrace": string(debug.Stack()),
			}).Error(err)

			http.Error(writer, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
			return
		}

		co.SetLoggerEntry(service.Log)
		co.Begin()
		defer rollbackCoreApplicationService(co, service.Log)

		context.Set(request, constants.CONTEXT_CORE, co)
		handler.ServeHTTP(writer, request)
	})
}

func rollbackCoreApplicationService(co *models.Core, log *logrus.Entry) {
	err := co.Rollback()
	if err != nil {
		log.WithFields(logrus.Fields{
			"event":      "phlapi::rollbackCoreApplicationService - Failed when trying to rollback Core",
			"stackTrace": string(debug.Stack()),
		}).Error(err)
	}
}

func InitMiddleware(database *sql.DB, log *logrus.Logger) Middleware {
	var middleware Middleware
	middleware.Database = database
	middleware.Log = log
	return middleware
}
