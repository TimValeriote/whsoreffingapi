package core

import (
	"database/sql"
	"runtime/debug"

	"github.com/sirupsen/logrus"
	"whosreffing-api/core/models"
)

type leagueStore struct {
	database *models.CoreDatabase
	log      *logrus.Entry
}

func LeagueStoreSetup(db *models.CoreDatabase, log *logrus.Entry) *leagueStore {
	return &leagueStore{
		database: db,
		log:      log,
	}
}

func (store *leagueStore) GetAllLeagues() ([]models.League, error) {
	sql := `SELECT id, name, level_id, code FROM league`
	query, err := store.database.Tx.Prepare(sql)
	if err != nil {
		store.log.WithFields(logrus.Fields{
			"event":      "activeGamesStore::GetAllLeagues - Failed to prepare GetAllLeagues SELECT query.",
			"stackTrace": string(debug.Stack()),
		}).Error(err)
		return nil, err
	}
	defer query.Close()

	ret, err := getLeaguesFromQuery(query)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (store *leagueStore) GetLeagueById(leagueId int) (models.League, error) {
	var league models.League

	sql := `SELECT id, name, level_id, code FROM league WHERE id = ?`

	err := store.database.Tx.QueryRow(sql, leagueId).Scan(
		&league.Id,
		&league.Name,
		&league.LevelId,
		&league.Code,
	)
	if err != nil {
		store.log.WithFields(logrus.Fields{
			"event":      "arenaStore::GetLeagueById - Failed to execute GetLeagueById SELECT query.",
			"stackTrace": string(debug.Stack()),
		}).Error(err)
		return league, err
	}

	return league, nil
}

func getLeaguesFromQuery(query *sql.Stmt) ([]models.League, error) {
	rows, err := query.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	leagues := make([]models.League, 0)
	for rows.Next() {
		var league models.League
		err = rows.Scan(
			&league.Id,
			&league.Name,
			&league.LevelId,
			&league.Code,
		)
		if err != nil {
			return nil, err
		}

		leagues = append(leagues, league)
	}

	if rows.Err() != nil {
		return nil, err
	}
	return leagues, nil
}

func (store *leagueStore) GetLeagueStandingsByLeagueId(leagueId int) ([]models.LeagueStanding, error) {
	sql := `SELECT 
				o.id AS official_id, 
				o.first_name, 
				o.last_name, 
				COUNT(og.game_id) AS games_worked 
			FROM official o 
				JOIN officials_games og ON o.id = og.official_id 
				JOIN game g ON og.game_id = g.id 
				JOIN league l ON g.league_id = l.id 
			WHERE l.id = ? 
				GROUP BY o.id, o.first_name, o.last_name 
				ORDER BY games_worked DESC`

	query, err := store.database.Tx.Prepare(sql)
	if err != nil {
		store.log.WithFields(logrus.Fields{
			"event":      "activeGamesStore::GetLeagueStandingsByLeagueId - Failed to prepare GetLeagueStandingsByLeagueId SELECT query.",
			"stackTrace": string(debug.Stack()),
		}).Error(err)
		return nil, err
	}
	defer query.Close()

	ret, err := getLeagueStandingsFromQuery(query, leagueId)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func getLeagueStandingsFromQuery(query *sql.Stmt, leagueId int) ([]models.LeagueStanding, error) {
	rows, err := query.Query(leagueId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	standings := make([]models.LeagueStanding, 0)
	for rows.Next() {
		var indiv models.LeagueStanding
		err = rows.Scan(
			&indiv.OfficialId,
			&indiv.FirstName,
			&indiv.LastName,
			&indiv.GameCount,
		)
		if err != nil {
			return nil, err
		}

		standings = append(standings, indiv)
	}

	if rows.Err() != nil {
		return nil, err
	}
	return standings, nil
}
