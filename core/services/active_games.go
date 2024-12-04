package core

import (
	"database/sql"
	"runtime/debug"

	"github.com/sirupsen/logrus"
	"whosreffing-api/core/models"
)

type activeGamesStore struct {
	database *models.CoreDatabase
	log      *logrus.Entry
}

func ActiveGamesStoreSetup(db *models.CoreDatabase, log *logrus.Entry) *activeGamesStore {
	return &activeGamesStore{
		database: db,
		log:      log,
	}
}

func (store *activeGamesStore) GetActiveGames() ([]models.ActiveGame, error) {
	sql := `SELECT id, game_id, ls_game_id FROM active_games`
	query, err := store.database.Tx.Prepare(sql)
	if err != nil {
		store.log.WithFields(logrus.Fields{
			"event":      "activeGamesStore::GetActiveGames - Failed to prepare GetActiveGames SELECT query.",
			"stackTrace": string(debug.Stack()),
		}).Error(err)
		return nil, err
	}
	defer query.Close()

	ret, err := getActiveGamesFromQuery(query)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func getActiveGamesFromQuery(query *sql.Stmt) ([]models.ActiveGame, error) {
	rows, err := query.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	activeGames := make([]models.ActiveGame, 0)
	for rows.Next() {
		var activeGame models.ActiveGame
		err = rows.Scan(
			&activeGame.Id,
			&activeGame.GameId,
			&activeGame.LsGameId,
		)
		if err != nil {
			return nil, err
		}

		activeGames = append(activeGames, activeGame)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return activeGames, nil
}

func (store *activeGamesStore) GetActiveGameDetails() ([]models.ActiveGameDetails, error) {
	sql := `SELECT
			    game_id,
			    active_games.ls_game_id,
			    league.name,
			    home_team.shortname AS home_team,
			    visitor_team.shortname AS visitor_team,
			    arena.name AS arena,
			    DATE_FORMAT(game.time, '%l:%i %p') AS game_time,
			    IF(game.referee_one_id = 0, 'No Official', CONCAT(referee_one.first_name, ' ', referee_one.last_name)) AS ref1,
			    IF(game.referee_two_id = 0, 'No Official', CONCAT(referee_two.first_name, ' ', referee_two.last_name)) AS ref2,
			    IF(game.linesman_one_id = 0, 'No Official', CONCAT(linesman_one.first_name, ' ', linesman_one.last_name)) AS ref3,
			    IF(game.linesman_two_id = 0, 'No Official', CONCAT(linesman_two.first_name, ' ', linesman_two.last_name)) AS ref4
			FROM active_games
			    JOIN game ON active_games.ls_game_id = game.ls_game_id
			    JOIN team AS home_team ON game.home_team_id = home_team.id
			    JOIN team AS visitor_team ON game.visiting_team_id = visitor_team.id
			    JOIN arena ON game.arena_id = arena.id
			    LEFT JOIN official AS referee_one ON game.referee_one_id = referee_one.id
			    LEFT JOIN official AS referee_two ON game.referee_two_id = referee_two.id
			    LEFT JOIN official AS linesman_one ON game.linesman_one_id = linesman_one.id
			    LEFT JOIN official AS linesman_two ON game.linesman_two_id = linesman_two.id
			    JOIN league ON home_team.league_id = league.id`
	query, err := store.database.Tx.Prepare(sql)
	if err != nil {
		store.log.WithFields(logrus.Fields{
			"event":      "activeGamesStore::GetActiveGameDetails - Failed to prepare GetActiveGameDetails SELECT query.",
			"stackTrace": string(debug.Stack()),
		}).Error(err)
		return nil, err
	}
	defer query.Close()

	ret, err := getActiveGamesDetailsFromQuery(query)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func getActiveGamesDetailsFromQuery(query *sql.Stmt) ([]models.ActiveGameDetails, error) {
	rows, err := query.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	activeGames := make([]models.ActiveGameDetails, 0)
	for rows.Next() {
		var activeGame models.ActiveGameDetails

		err = rows.Scan(
			&activeGame.GameId,
			&activeGame.LsGameId,
			&activeGame.LeagueName,
			&activeGame.HomeTeam,
			&activeGame.VisitorTeam,
			&activeGame.Arena,
			&activeGame.GameTime,
			&activeGame.RefereeOne,
			&activeGame.RefereeTwo,
			&activeGame.LinesmanOne,
			&activeGame.LinesmanTwo,
		)
		if err != nil {
			return nil, err
		}

		if activeGame.RefereeOne == "No Official" && activeGame.RefereeTwo == "No Official" && activeGame.LinesmanOne == "No Official" && activeGame.LinesmanTwo == "No Official" {
			activeGame.HasNoOfficials = true
		} else {
			activeGame.HasNoOfficials = false
		}

		activeGames = append(activeGames, activeGame)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return activeGames, nil
}
