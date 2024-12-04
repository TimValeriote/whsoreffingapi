package implementation

import (
	"whosreffing-api/core/models"
	services "whosreffing-api/core/services"
)

type GameTypeStore struct {
	*models.Core
}

func (store *GameTypeStore) GetGameTypeById(gameTypeId int) (models.GameType, error) {
	return services.GameTypeStoreSetup(store.Database, store.Log).GetGameTypeById(gameTypeId)
}
