package implementation

import (
	"whosreffing-api/core/models"
	services "whosreffing-api/core/services"
)

type ActiveGamesStore struct {
	*models.Core
}

func (store *ActiveGamesStore) GetActiveGames() ([]models.ActiveGame, error) {
	return services.ActiveGamesStoreSetup(store.Database, store.Log).GetActiveGames()
}

func (store *ActiveGamesStore) GetActiveGameDetails() ([]models.ActiveGameDetails, error) {
	return services.ActiveGamesStoreSetup(store.Database, store.Log).GetActiveGameDetails()
}
