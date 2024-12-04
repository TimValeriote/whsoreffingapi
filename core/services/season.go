package core

import (
	"database/sql"
	"runtime/debug"

	"github.com/sirupsen/logrus"
	"whosreffing-api/core/models"
)

type seasonStore struct {
	database *models.CoreDatabase
	log      *logrus.Entry
}

func SeasonStoreSetup(db *models.CoreDatabase, log *logrus.Entry) *seasonStore {
	return &seasonStore{
		database: db,
		log:      log,
	}
}

func (store *seasonStore) GetAllActiveSeasons() ([]models.Season, error) {
	sql := `SELECT id, season_name, season_type_id, league_id, ls_season_id, active FROM season WHERE active = 1`
	query, err := store.database.Tx.Prepare(sql)
	if err != nil {
		store.log.WithFields(logrus.Fields{
			"event":      "activeGamesStore::GetAllActiveSeasons - Failed to prepare GetAllActiveSeasons SELECT query.",
			"stackTrace": string(debug.Stack()),
		}).Error(err)
		return nil, err
	}
	defer query.Close()

	ret, err := getSeasonsFromQuery(query)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (store *seasonStore) GetSeasonById(seasonId int) (models.Season, error) {
	var season models.Season

	sql := `SELECT id, season_name, season_type_id, league_id, ls_season_id, active FROM season WHERE id = ?`

	err := store.database.Tx.QueryRow(sql, seasonId).Scan(
		&season.Id,
		&season.SeasonName,
		&season.SeasonTypeId,
		&season.LeagueId,
		&season.LsSeasonId,
		&season.Active,
	)
	if err != nil {
		store.log.WithFields(logrus.Fields{
			"event":      "arenaStore::GetSeasonById - Failed to execute GetSeasonById SELECT query.",
			"stackTrace": string(debug.Stack()),
		}).Error(err)
		return season, err
	}

	return season, nil
}

func getSeasonsFromQuery(query *sql.Stmt) ([]models.Season, error) {
	rows, err := query.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	seasons := make([]models.Season, 0)
	for rows.Next() {
		var season models.Season
		err = rows.Scan(
			&season.Id,
			&season.SeasonName,
			&season.SeasonTypeId,
			&season.LeagueId,
			&season.LsSeasonId,
			&season.Active,
		)
		if err != nil {
			return nil, err
		}

		seasons = append(seasons, season)
	}

	if rows.Err() != nil {
		return nil, err
	}
	return seasons, nil
}
