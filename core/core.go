package core

import (
	"context"
	"database/sql"

	"whosreffing-api/core/implementation"
	"whosreffing-api/core/models"
)

func CreateCore(db *sql.DB) (*models.Core, error) {
	return CreateBACoreContext(context.Background(), db)
}

func CreateBACoreContext(ctx context.Context, db *sql.DB) (*models.Core, error) {
	var coreDB models.CoreDatabase

	core := &models.Core{
		Ctx: ctx,
	}

	coreDB.DB = db
	core.Database = &coreDB

	core.ActiveGamesService = &implementation.ActiveGamesStore{core}
	core.ArenaService = &implementation.ArenaStore{core}
	core.GameService = &implementation.GameStore{core}
	core.GameDetailsService = &implementation.GameDetailsStore{core}
	core.GameTypeService = &implementation.GameTypeStore{core}
	core.KeysService = &implementation.KeysStore{core}
	core.LeagueService = &implementation.LeagueStore{core}
	core.OfficialService = &implementation.OfficialStore{core}
	core.SeasonService = &implementation.SeasonStore{core}
	core.StatusService = &implementation.StatusStore{core}
	core.TeamService = &implementation.TeamStore{core}
	core.OfficialGameService = &implementation.OfficialGameStore{core}

	return core, nil
}
