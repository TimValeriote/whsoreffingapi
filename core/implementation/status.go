package implementation

import (
	"whosreffing-api/core/models"
	services "whosreffing-api/core/services"
)

type StatusStore struct {
	*models.Core
}

func (store *StatusStore) GetStatusById(statusId int) (models.Status, error) {
	return services.StatusStoreSetup(store.Database, store.Log).GetStatusById(statusId)
}
