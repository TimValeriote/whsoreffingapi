package implementation

import (
	"whosreffing-api/core/models"
	services "whosreffing-api/core/services"
)

type OfficialStore struct {
	*models.Core
}

func (store *OfficialStore) GetOfficialById(officialId int) (models.Official, error) {
	return services.OfficialStoreSetup(store.Database, store.Log).GetOfficialById(officialId)
}

func (store *OfficialStore) GetOfficialInfoByOfficialId(officialId int) (models.OfficialInfo, error) {
	return services.OfficialStoreSetup(store.Database, store.Log).GetOfficialInfoByOfficialId(officialId)
}

func (store *OfficialStore) GetOfficialsBySearchTerm(searchTerm string) ([]models.OfficialSearch, error) {
	return services.OfficialStoreSetup(store.Database, store.Log).GetOfficialsBySearchTerm(searchTerm)
}
