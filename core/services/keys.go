package core

import (
	"database/sql"
	"runtime/debug"

	"github.com/sirupsen/logrus"
	"whosreffing-api/core/models"
)

type keysStore struct {
	database *models.CoreDatabase
	log      *logrus.Entry
}

func KeysStoreSetup(db *models.CoreDatabase, log *logrus.Entry) *keysStore {
	return &keysStore{
		database: db,
		log:      log,
	}
}

func (store *keysStore) GetAllKeys() ([]models.Key, error) {
	sql := `SELECT id, description, league_id FROM keys`
	query, err := store.database.Tx.Prepare(sql)
	if err != nil {
		store.log.WithFields(logrus.Fields{
			"event":      "activeGamesStore::GetAllKeys - Failed to prepare GetAllKeys SELECT query.",
			"stackTrace": string(debug.Stack()),
		}).Error(err)
		return nil, err
	}
	defer query.Close()

	ret, err := getKeysFromQuery(query)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (store *keysStore) GetKeyByLeagueId(leagueId int) (models.Key, error) {
	var key models.Key

	sql := `SELECT id, description, league_id FROM keys WHERE league_id = ?`

	err := store.database.Tx.QueryRow(sql, leagueId).Scan(
		&key.Id,
		&key.Description,
		&key.LeagueId,
	)
	if err != nil {
		store.log.WithFields(logrus.Fields{
			"event":      "arenaStore::GetKeyByLeagueId - Failed to execute GetKeyByLeagueId SELECT query.",
			"stackTrace": string(debug.Stack()),
		}).Error(err)
		return key, err
	}

	return key, nil
}

func getKeysFromQuery(query *sql.Stmt) ([]models.Key, error) {
	rows, err := query.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	keys := make([]models.Key, 0)
	for rows.Next() {
		var key models.Key
		err = rows.Scan(
			&key.Id,
			&key.Description,
			&key.LeagueId,
		)
		if err != nil {
			return nil, err
		}

		keys = append(keys, key)
	}

	if rows.Err() != nil {
		return nil, err
	}
	return keys, nil
}
