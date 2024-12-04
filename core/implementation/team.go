package implementation

import (
	"whosreffing-api/core/models"
	services "whosreffing-api/core/services"
)

type TeamStore struct {
	*models.Core
}

func (store *TeamStore) GetTeamById(teamId int) (models.Team, error) {
	return services.TeamStoreSetup(store.Database, store.Log).GetTeamById(teamId)
}
