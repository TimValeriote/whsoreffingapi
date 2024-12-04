package core

import (
	"runtime/debug"

	"github.com/sirupsen/logrus"
	"whosreffing-api/core/models"
)

type teamStore struct {
	database *models.CoreDatabase
	log      *logrus.Entry
}

func TeamStoreSetup(db *models.CoreDatabase, log *logrus.Entry) *teamStore {
	return &teamStore{
		database: db,
		log:      log,
	}
}

func (store *teamStore) GetTeamById(teamId int) (models.Team, error) {
	var team models.Team

	sql := `SELECT id, name, nickname, league_id, code FROM team WHERE id = ?`

	err := store.database.Tx.QueryRow(sql, teamId).Scan(
		&team.Id,
		&team.Name,
		&team.Code,
		&team.NickName,
		&team.LeagueId,
	)
	if err != nil {
		store.log.WithFields(logrus.Fields{
			"event":      "arenaStore::GetTeamById - Failed to execute GetTeamById SELECT query.",
			"stackTrace": string(debug.Stack()),
		}).Error(err)
		return team, err
	}

	return team, nil
}
