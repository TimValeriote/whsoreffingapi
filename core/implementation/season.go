package implementation

import (
	"whosreffing-api/core/models"
	services "whosreffing-api/core/services"
)

type SeasonStore struct {
	*models.Core
}

func (store *SeasonStore) GetAllActiveSeasons() ([]models.Season, error) {
	return services.SeasonStoreSetup(store.Database, store.Log).GetAllActiveSeasons()
}

func (store *SeasonStore) GetSeasonById(seasonId int) (models.Season, error) {
	return services.SeasonStoreSetup(store.Database, store.Log).GetSeasonById(seasonId)
}
