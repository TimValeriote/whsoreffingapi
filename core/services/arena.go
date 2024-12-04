package core

import (
	"runtime/debug"

	"github.com/sirupsen/logrus"
	"whosreffing-api/core/models"
)

type arenaStore struct {
	database *models.CoreDatabase
	log      *logrus.Entry
}

func ArenaStoreSetup(db *models.CoreDatabase, log *logrus.Entry) *arenaStore {
	return &arenaStore{
		database: db,
		log:      log,
	}
}

func (store *arenaStore) GetArenaById(arenaId int) (models.Arena, error) {
	var arena models.Arena

	sql := `SELECT id, name FROM arena WHERE id = ?`

	err := store.database.Tx.QueryRow(sql, arenaId).Scan(&arena.Id, &arena.Name)
	if err != nil {
		store.log.WithFields(logrus.Fields{
			"event":      "arenaStore::GetArenaById - Failed to execute GetArenaById SELECT query.",
			"stackTrace": string(debug.Stack()),
		}).Error(err)
		return arena, err
	}

	return arena, nil
}
