package implementation

import (
	"whosreffing-api/core/models"
	services "whosreffing-api/core/services"
)

type LeagueStore struct {
	*models.Core
}

func (store *LeagueStore) GetAllLeagues() ([]models.League, error) {
	return services.LeagueStoreSetup(store.Database, store.Log).GetAllLeagues()
}

func (store *LeagueStore) GetLeagueById(leagueId int) (models.League, error) {
	return services.LeagueStoreSetup(store.Database, store.Log).GetLeagueById(leagueId)
}

func (store *LeagueStore) GetLeagueStandingsByLeagueId(leagueId int) ([]models.LeagueStanding, error) {
	return services.LeagueStoreSetup(store.Database, store.Log).GetLeagueStandingsByLeagueId(leagueId)
}
