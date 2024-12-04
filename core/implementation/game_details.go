package implementation

import (
	"whosreffing-api/core/models"
	services "whosreffing-api/core/services"
)

type GameDetailsStore struct {
	*models.Core
}

func (store *GameDetailsStore) GetGameDetailsByGameId(gameId int) (models.GameDetails, error) {
	return services.GameDetailsStoreSetup(store.Database, store.Log).GetGameDetailsByGameId(gameId)
}

func (store *GameDetailsStore) GetOfficialGameDetailsByOfficialID(officialId int) ([]models.GameDetails, error) {
	return services.GameDetailsStoreSetup(store.Database, store.Log).GetOfficialGameDetailsByOfficialID(officialId)
}
