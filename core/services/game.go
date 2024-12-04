package core

import (
	"database/sql"
	"runtime/debug"

	"github.com/sirupsen/logrus"
	"whosreffing-api/core/models"
)

type gameStore struct {
	database *models.CoreDatabase
	log      *logrus.Entry
}

func GameStoreSetup(db *models.CoreDatabase, log *logrus.Entry) *gameStore {
	return &gameStore{
		database: db,
		log:      log,
	}
}

func (store *gameStore) GetGameById(gameId int) (models.Game, error) {
	var game models.Game

	sql := `SELECT 
		id,
		ls_game_id, 
		league_id,
		date,
		time,
		home_team_id,
		visiting_team_id,
		arena_id,
		status_id,
		referee_one_id,
		referee_two_id,
		linesman_one_id,
		linesman_two_id,
		season_id
	FROM game WHERE id = ?`

	err := store.database.Tx.QueryRow(sql, gameId).Scan(
		&game.Id,
		&game.LsGameID,
		&game.LeagueId,
		&game.Date,
		&game.Time,
		&game.HomeTeamId,
		&game.VisitingTeamId,
		&game.ArenaId,
		&game.StatusId,
		&game.RefereeOneId,
		&game.RefereeTwoId,
		&game.LinesmanOneId,
		&game.LinesmanTwoId,
		&game.SeasonId,
	)
	if err != nil {
		store.log.WithFields(logrus.Fields{
			"event":      "arenaStore::GetGameById - Failed to execute GetGameById SELECT query.",
			"stackTrace": string(debug.Stack()),
		}).Error(err)
		return game, err
	}

	return game, nil

}

func (store *gameStore) GetGamesCalendar() ([]models.GamesCalendar, error) {
	sql := `SELECT DATE(date) AS game_date, COUNT(*) AS total_games
			FROM game
			GROUP BY game_date
			ORDER BY game_date`
	query, err := store.database.Tx.Prepare(sql)
	if err != nil {
		store.log.WithFields(logrus.Fields{
			"event":      "activeGamesStore::GetGamesCalendar - Failed to prepare GetGamesCalendar SELECT query.",
			"stackTrace": string(debug.Stack()),
		}).Error(err)
		return nil, err
	}
	defer query.Close()

	ret, err := getGamesCalendarByQuery(query)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func getGamesCalendarByQuery(query *sql.Stmt) ([]models.GamesCalendar, error) {
	rows, err := query.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	gamesCalendar := make([]models.GamesCalendar, 0)
	for rows.Next() {
		var game models.GamesCalendar
		err = rows.Scan(
			&game.Date,
			&game.GameCount,
		)
		if err != nil {
			return nil, err
		}

		gamesCalendar = append(gamesCalendar, game)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return gamesCalendar, nil
}
