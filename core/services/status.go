package core

import (
	"runtime/debug"

	"github.com/sirupsen/logrus"
	"whosreffing-api/core/models"
)

type statusStore struct {
	database *models.CoreDatabase
	log      *logrus.Entry
}

func StatusStoreSetup(db *models.CoreDatabase, log *logrus.Entry) *statusStore {
	return &statusStore{
		database: db,
		log:      log,
	}
}

func (store *statusStore) GetStatusById(statusId int) (models.Status, error) {
	var status models.Status

	sql := `SELECT SELECT id, description FROM status_trans WHERE id = ?`

	err := store.database.Tx.QueryRow(sql, statusId).Scan(
		&status.Id,
		&status.Description,
	)
	if err != nil {
		store.log.WithFields(logrus.Fields{
			"event":      "arenaStore::GetStatusById - Failed to execute GetStatusById SELECT query.",
			"stackTrace": string(debug.Stack()),
		}).Error(err)
		return status, err
	}

	return status, nil
}
