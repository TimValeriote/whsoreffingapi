package main

import (
	"database/sql"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/sirupsen/logrus"
	"whosreffing-api/apibuilder"
	"whosreffing-api/controllers"
)

func SetupRouting(router *httprouter.Router, database *sql.DB, log *logrus.Logger) {
	middleware := InitMiddleware(database, log)
	api := apibuilder.NewApi("/whosreffing", router, alice.New(context.ClearHandler, middleware.CORSMiddleware), alice.New(context.ClearHandler, middleware.OptionsMiddleware), log)

	api.Routes = []apibuilder.Route{
		//This is the index route, really only for testing and to confirm the API is running on your local enviro
		{"GET", "/index", controllers.IndexController{Log: log}.Index, alice.New()},

		{"GET", "/activegames", controllers.ActiveGamesController{Log: log}.GetActiveGames, alice.New(middleware.CoreMasterCoreMiddleware, middleware.CoreApplicationServiceMiddleware)},
		{"GET", "/activegamesdetails", controllers.ActiveGamesController{Log: log}.GetActiveGamesDetails, alice.New(middleware.CoreMasterCoreMiddleware, middleware.CoreApplicationServiceMiddleware)},

		{"GET", "/arena/:arenaId", controllers.ArenaController{Log: log}.GetArenaById, alice.New(middleware.CoreMasterCoreMiddleware, middleware.CoreApplicationServiceMiddleware)},

		{"GET", "/game/:gameId", controllers.GameController{Log: log}.GetGameById, alice.New(middleware.CoreMasterCoreMiddleware, middleware.CoreApplicationServiceMiddleware)},
		{"GET", "/gamescalendar", controllers.GameController{Log: log}.GetGamesCalendar, alice.New(middleware.CoreMasterCoreMiddleware, middleware.CoreApplicationServiceMiddleware)},

		{"GET", "/gamedetail/:gameId", controllers.GameDetailsController{Log: log}.GetGameDetailsByGameId, alice.New(middleware.CoreMasterCoreMiddleware, middleware.CoreApplicationServiceMiddleware)},
		{"GET", "/gamedetails/:officialId", controllers.GameDetailsController{Log: log}.GetOfficialGamesByOfficialId, alice.New(middleware.CoreMasterCoreMiddleware, middleware.CoreApplicationServiceMiddleware)},

		{"GET", "/official/:officialId", controllers.OfficialController{Log: log}.GetOfficialById, alice.New(middleware.CoreMasterCoreMiddleware, middleware.CoreApplicationServiceMiddleware)},
		{"GET", "/officialgames/:officialId", controllers.OfficialGameController{Log: log}.GetOfficialGamesByOfficialId, alice.New(middleware.CoreMasterCoreMiddleware, middleware.CoreApplicationServiceMiddleware)},
		{"GET", "/officialdetails/:officialId", controllers.OfficialController{Log: log}.GetOfficialDetailsById, alice.New(middleware.CoreMasterCoreMiddleware, middleware.CoreApplicationServiceMiddleware)},
		{"GET", "/officialssearch", controllers.OfficialController{Log: log}.GetOfficialsBySearchTerm, alice.New(middleware.CoreMasterCoreMiddleware, middleware.CoreApplicationServiceMiddleware)},

		{"GET", "/leagues", controllers.LeagueController{Log: log}.GetAllLeagues, alice.New(middleware.CoreMasterCoreMiddleware, middleware.CoreApplicationServiceMiddleware)},
		{"GET", "/league/:leagueId", controllers.LeagueController{Log: log}.GetLeagueById, alice.New(middleware.CoreMasterCoreMiddleware, middleware.CoreApplicationServiceMiddleware)},
		{"GET", "/leaguestandings/:leagueId", controllers.LeagueController{Log: log}.GetLeagueStandingsByLeagueId, alice.New(middleware.CoreMasterCoreMiddleware, middleware.CoreApplicationServiceMiddleware)},
	}

	api.Finalize()
}
