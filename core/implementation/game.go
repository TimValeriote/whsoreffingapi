package implementation

import (
	"whosreffing-api/core/models"
	services "whosreffing-api/core/services"
)

type GameStore struct {
	*models.Core
}

func (store *GameStore) GetGameById(gameId int) (models.Game, error) {
	return services.GameStoreSetup(store.Database, store.Log).GetGameById(gameId)
}

func (store *GameStore) GetGamesCalendar() ([]models.GamesCalendar, error) {
	return services.GameStoreSetup(store.Database, store.Log).GetGamesCalendar()
}
