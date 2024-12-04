package implementation

import (
	"whosreffing-api/core/models"
	services "whosreffing-api/core/services"
)

type ArenaStore struct {
	*models.Core
}

func (store *ArenaStore) GetArenaById(arenaId int) (models.Arena, error) {
	return services.ArenaStoreSetup(store.Database, store.Log).GetArenaById(arenaId)
}
