package core

import (
	"runtime/debug"

	"github.com/sirupsen/logrus"
	"whosreffing-api/core/models"
)

type officialGameStore struct {
	database *models.CoreDatabase
	log      *logrus.Entry
}

func OfficialGameStoreSetup(db *models.CoreDatabase, log *logrus.Entry) *officialGameStore {
	return &officialGameStore{
		database: db,
		log:      log,
	}
}

func (store *officialGameStore) GetOfficialGamesByOfficialId(officialId int) ([]models.OfficialGame, error) {
	sql := `SELECT id, game_id, official_id FROM officials_games WHERE official_id = ?`
	query, err := store.database.Tx.Prepare(sql)
	if err != nil {
		store.log.WithFields(logrus.Fields{
			"event":      "officialGameStore::GetOfficialGamesByOfficialId - Failed to prepare GetOfficialGamesByOfficialId SELECT query.",
			"stackTrace": string(debug.Stack()),
		}).Error(err)
		return nil, err
	}
	defer query.Close()

	rows, err := query.Query(officialId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	games := make([]models.OfficialGame, 0)
	for rows.Next() {
		var game models.OfficialGame
		err = rows.Scan(
			&game.Id,
			&game.GameId,
			&game.OfficialId,
		)
		if err != nil {
			return nil, err
		}

		games = append(games, game)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return games, nil
}
