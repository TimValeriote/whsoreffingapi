package implementation

import (
	"whosreffing-api/core/models"
	services "whosreffing-api/core/services"
)

type KeysStore struct {
	*models.Core
}

func (store *KeysStore) GetAllKeys() ([]models.Key, error) {
	return services.KeysStoreSetup(store.Database, store.Log).GetAllKeys()
}

func (store *KeysStore) GetKeyByLeagueId(leagueId int) (models.Key, error) {
	return services.KeysStoreSetup(store.Database, store.Log).GetKeyByLeagueId(leagueId)
}
