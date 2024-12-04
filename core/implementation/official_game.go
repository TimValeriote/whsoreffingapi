package implementation

import (
	"whosreffing-api/core/models"
	services "whosreffing-api/core/services"
)

type OfficialGameStore struct {
	*models.Core
}

func (store *OfficialGameStore) GetOfficialGamesByOfficialId(officialId int) ([]models.OfficialGame, error) {
	return services.OfficialGameStoreSetup(store.Database, store.Log).GetOfficialGamesByOfficialId(officialId)
}
