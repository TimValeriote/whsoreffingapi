package core

import (
	"runtime/debug"

	"github.com/sirupsen/logrus"
	"whosreffing-api/core/models"
)

type gameDetailsStore struct {
	database *models.CoreDatabase
	log      *logrus.Entry
}

func GameDetailsStoreSetup(db *models.CoreDatabase, log *logrus.Entry) *gameDetailsStore {
	return &gameDetailsStore{
		database: db,
		log:      log,
	}
}

func (store *gameDetailsStore) GetGameDetailsByGameId(gameId int) (models.GameDetails, error) {
	var gameDetails models.GameDetails

	sql := `SELECT 
		id,
		game_id, 
		game_length,
		home_goals,
		visitor_goals,
		overtime,
		shootout,
		home_powerplays,
		visitor_powerplays,
		home_pims,
		visitor_pims,
		home_faceoffs,
		visitor_faceoffs,
		home_faceoffs_won,
		visitor_faceoffs_won
	FROM game_details WHERE game_id = ?`

	err := store.database.Tx.QueryRow(sql, gameId).Scan(
		&gameDetails.Id,
		&gameDetails.GameId,
		&gameDetails.GameLength,
		&gameDetails.HomeGoals,
		&gameDetails.VisitingGoals,
		&gameDetails.OverTime,
		&gameDetails.Shootout,
		&gameDetails.HomePowerPlays,
		&gameDetails.VisitorPowerPlays,
		&gameDetails.HomePims,
		&gameDetails.VisitorPims,
		&gameDetails.HomeFaceOffs,
		&gameDetails.VisitorFaceOffs,
		&gameDetails.HomeFaceOffsWon,
		&gameDetails.VisitorFaceOffsWon,
	)
	if err != nil {
		store.log.WithFields(logrus.Fields{
			"event":      "arenaStore::GetGameDetailsByGameId - Failed to execute GetGameDetailsByGameId SELECT query.",
			"stackTrace": string(debug.Stack()),
		}).Error(err)
		return gameDetails, err
	}

	return gameDetails, nil
}

func (store *gameDetailsStore) GetOfficialGameDetailsByOfficialID(officialId int) ([]models.GameDetails, error) {
	sql := `SELECT 
		gd.id,
		gd.game_id, 
		gd.game_length,
		gd.home_goals,
		gd.visitor_goals,
		gd.overtime,
		gd.shootout,
		gd.home_powerplays,
		gd.visitor_powerplays,
		gd.home_pims,
		gd.visitor_pims,
		gd.home_faceoffs,
		gd.visitor_faceoffs,
		gd.home_faceoffs_won,
		gd.visitor_faceoffs_won
	FROM game_details AS gd 
	INNER JOIN game AS g ON gd.game_id = g.id
	WHERE ? IN (referee_one_id, referee_two_id, linesman_one_id, linesman_two_id)`

	query, err := store.database.Tx.Prepare(sql)
	if err != nil {
		store.log.WithFields(logrus.Fields{
			"event":      "officialGameStore::GetOfficialGameDetailsByOfficialID - Failed to prepare GetOfficialGameDetailsByOfficialID SELECT query.",
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

	games := make([]models.GameDetails, 0)
	for rows.Next() {
		var gameDetails models.GameDetails
		err = rows.Scan(
			&gameDetails.Id,
			&gameDetails.GameId,
			&gameDetails.GameLength,
			&gameDetails.HomeGoals,
			&gameDetails.VisitingGoals,
			&gameDetails.OverTime,
			&gameDetails.Shootout,
			&gameDetails.HomePowerPlays,
			&gameDetails.VisitorPowerPlays,
			&gameDetails.HomePims,
			&gameDetails.VisitorPims,
			&gameDetails.HomeFaceOffs,
			&gameDetails.VisitorFaceOffs,
			&gameDetails.HomeFaceOffsWon,
			&gameDetails.VisitorFaceOffsWon,
		)
		if err != nil {
			return nil, err
		}

		games = append(games, gameDetails)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return games, nil

}
