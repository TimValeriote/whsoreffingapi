package core

import (
	"runtime/debug"

	"github.com/sirupsen/logrus"
	"whosreffing-api/core/models"
)

type gameTypeStore struct {
	database *models.CoreDatabase
	log      *logrus.Entry
}

func GameTypeStoreSetup(db *models.CoreDatabase, log *logrus.Entry) *gameTypeStore {
	return &gameTypeStore{
		database: db,
		log:      log,
	}
}

func (store *gameTypeStore) GetGameTypeById(gameTypeId int) (models.GameType, error) {
	var gameType models.GameType

	sql := `SELECT id, description FROM game_type_trans WHERE id = ?`

	err := store.database.Tx.QueryRow(sql, gameTypeId).Scan(
		&gameType.Id,
		&gameType.Description,
	)
	if err != nil {
		store.log.WithFields(logrus.Fields{
			"event":      "arenaStore::GetGameTypeById - Failed to execute GetGameTypeById SELECT query.",
			"stackTrace": string(debug.Stack()),
		}).Error(err)
		return gameType, err
	}

	return gameType, nil
}
