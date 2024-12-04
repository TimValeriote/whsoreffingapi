package core

import (
	"database/sql"
	"fmt"
	"math"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"whosreffing-api/core/models"
)

type officialStore struct {
	database *models.CoreDatabase
	log      *logrus.Entry
}

func OfficialStoreSetup(db *models.CoreDatabase, log *logrus.Entry) *officialStore {
	return &officialStore{
		database: db,
		log:      log,
	}
}

func (store *officialStore) GetOfficialById(officialId int) (models.Official, error) {
	var official models.Official

	sql := `SELECT id, first_name, last_name, jersey_number FROM official WHERE id = ?`

	err := store.database.Tx.QueryRow(sql, officialId).Scan(
		&official.Id,
		&official.FirstName,
		&official.LastName,
		&official.JerseyNumber,
	)
	if err != nil {
		store.log.WithFields(logrus.Fields{
			"event":      "arenaStore::GetOfficialById - Failed to execute GetOfficialById SELECT query.",
			"stackTrace": string(debug.Stack()),
		}).Error(err)
		return official, err
	}

	return official, nil
}

func (store *officialStore) GetOfficialsBySearchTerm(searchTerm string) ([]models.OfficialSearch, error) {
	var searchResults []models.OfficialSearch

	if strings.Contains(searchTerm, " ") {
		nameParts := strings.Split(searchTerm, " ")

		if len(nameParts) == 2 {
			firstName := nameParts[0]
			lastName := nameParts[1]

			sql := `SELECT
						o.id,
					    o.first_name,
					    o.last_name,
					    o.jersey_number,
					    GROUP_CONCAT(DISTINCT l.name) AS leagues
					FROM
					    official o
					JOIN
					    officials_games og ON o.id = og.official_id
					JOIN
					    game g ON og.game_id = g.id
					JOIN
					    league l ON g.league_id = l.id
					WHERE
					    o.first_name LIKE ? AND o.last_name LIKE ?
					GROUP BY
					    o.id`

			query, err := store.database.Tx.Prepare(sql)
			if err != nil {
				store.log.WithFields(logrus.Fields{
					"event":      "activeGamesStore::GetOfficialsBySearchTerm - Failed to prepare GetOfficialsBySearchTerm SELECT query.",
					"stackTrace": string(debug.Stack()),
				}).Error(err)
				return nil, err
			}
			defer query.Close()

			lastNameTermForQuery := lastName + "%"

			results, err := getOfficialsSearchGamesFromQueryWithFirstName(query, firstName, lastNameTermForQuery)
			if err != nil {
				store.log.WithFields(logrus.Fields{
					"event":      "activeGamesStore::GetOfficialsBySearchTerm - Failed to execute SELECT query.",
					"stackTrace": string(debug.Stack()),
				}).Error(err)
				return nil, err
			}

			searchResults = append(searchResults, results...)

			return searchResults, nil
		}
	} else {

		var searchResults []models.OfficialSearch

		firstNameSql := `SELECT
						o.id,
					    o.first_name,
					    o.last_name,
					    o.jersey_number,
					    GROUP_CONCAT(DISTINCT l.name) AS leagues
					FROM
					    official o
					JOIN
					    officials_games og ON o.id = og.official_id
					JOIN
					    game g ON og.game_id = g.id
					JOIN
					    league l ON g.league_id = l.id
					WHERE first_name LIKE ?
					GROUP BY o.id`

		firstNameQuery, err := store.database.Tx.Prepare(firstNameSql)
		if err != nil {
			store.log.WithFields(logrus.Fields{
				"event":      "activeGamesStore::GetOfficialsBySearchTerm - Failed to prepare GetOfficialsBySearchTerm SELECT query.",
				"stackTrace": string(debug.Stack()),
			}).Error(err)
			return nil, err
		}
		defer firstNameQuery.Close()

		searchTermForQuery := searchTerm + "%"

		firstNameResults, err := getOfficialsSearchGamesFromQuery(firstNameQuery, searchTermForQuery)
		if err != nil {
			store.log.WithFields(logrus.Fields{
				"event":      "activeGamesStore::GetOfficialsBySearchTerm - Failed to execute SELECT query.",
				"stackTrace": string(debug.Stack()),
			}).Error(err)
			return nil, err
		}

		searchResults = append(searchResults, firstNameResults...)

		lastNameSql := `SELECT
						o.id,
					    o.first_name,
					    o.last_name,
					    o.jersey_number,
					    GROUP_CONCAT(DISTINCT l.name) AS leagues
					FROM
					    official o
					JOIN
					    officials_games og ON o.id = og.official_id
					JOIN
					    game g ON og.game_id = g.id
					JOIN
					    league l ON g.league_id = l.id
					WHERE last_name LIKE ?
					GROUP BY o.id`

		lastNameQuery, err := store.database.Tx.Prepare(lastNameSql)
		if err != nil {
			store.log.WithFields(logrus.Fields{
				"event":      "activeGamesStore::GetOfficialsBySearchTerm - Failed to prepare GetOfficialsBySearchTerm SELECT query.",
				"stackTrace": string(debug.Stack()),
			}).Error(err)
			return nil, err
		}
		defer lastNameQuery.Close()

		lastNameResults, err := getOfficialsSearchGamesFromQuery(lastNameQuery, searchTermForQuery)
		if err != nil {
			store.log.WithFields(logrus.Fields{
				"event":      "activeGamesStore::GetOfficialsBySearchTerm - Failed to execute SELECT query.",
				"stackTrace": string(debug.Stack()),
			}).Error(err)
			return nil, err
		}

		searchResults = append(searchResults, lastNameResults...)

		return searchResults, nil
	}

	return searchResults, nil
}

func getOfficialsSearchGamesFromQueryWithFirstName(query *sql.Stmt, firstName string, lastName string) ([]models.OfficialSearch, error) {
	rows, err := query.Query(firstName, lastName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	officials := make([]models.OfficialSearch, 0)
	for rows.Next() {
		var official models.OfficialSearch

		err = rows.Scan(
			&official.Id,
			&official.FirstName,
			&official.LastName,
			&official.JerseyNumber,
			&official.Leagues,
		)
		if err != nil {
			return nil, err
		}
		officials = append(officials, official)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return officials, nil
}

func getOfficialsSearchGamesFromQuery(query *sql.Stmt, searchTerm string) ([]models.OfficialSearch, error) {
	rows, err := query.Query(searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	officials := make([]models.OfficialSearch, 0)
	for rows.Next() {
		var official models.OfficialSearch

		err = rows.Scan(
			&official.Id,
			&official.FirstName,
			&official.LastName,
			&official.JerseyNumber,
			&official.Leagues,
		)
		if err != nil {
			return nil, err
		}
		officials = append(officials, official)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return officials, nil
}

func (store *officialStore) GetOfficialInfoByOfficialId(officialId int) (models.OfficialInfo, error) {
	var officialInfo models.OfficialInfo

	officialInfo.Id = officialId

	sql := `SELECT first_name, last_name, jersey_number FROM official WHERE id = ?`

	err := store.database.Tx.QueryRow(sql, officialId).Scan(
		&officialInfo.FirstName,
		&officialInfo.LastName,
		&officialInfo.JerseyNumber,
	)
	if err != nil {
		store.log.WithFields(logrus.Fields{
			"event":      "arenaStore::GetOfficialById - Failed to execute GetOfficialById SELECT query.",
			"stackTrace": string(debug.Stack()),
		}).Error(err)
		return officialInfo, err
	}

	gamesSql := `SELECT 
			g.id AS GameId, 
			g.ls_game_id AS LsGameId, 
			l.name AS LeagueName,
			l.id AS LeagueId,
			th.shortname AS HomeTeam, 
			tv.shortname AS VisitorTeam, 
			a.name AS Arena,
			DATE_FORMAT(g.time, '%l:%i %p') AS game_time,
			DATE_FORMAT(g.date, '%b %e, %Y') as game_date,
			CONCAT(ro.first_name, ' ', ro.last_name) AS RefereeOne,
			CONCAT(rt.first_name, ' ', rt.last_name) AS RefereeTwo,
			CONCAT(l1t.first_name, ' ', l1t.last_name) AS LinesmanOne,
			CONCAT(l2t.first_name, ' ', l2t.last_name) AS LinesmanTwo
		FROM officials_games og
			JOIN official o ON og.official_id = o.id
			JOIN game g ON og.game_id = g.id
			JOIN league l ON g.league_id = l.id
			JOIN team th ON g.home_team_id = th.id
			JOIN team tv ON g.visiting_team_id = tv.id
			JOIN arena a ON g.arena_id = a.id
			LEFT JOIN official ro ON g.referee_one_id = ro.id
			LEFT JOIN official rt ON g.referee_two_id = rt.id
			LEFT JOIN official l1t ON g.linesman_one_id = l1t.id
			LEFT JOIN official l2t ON g.linesman_two_id = l2t.id
		WHERE o.id = ?
		ORDER BY game_date ASC`

	query, err := store.database.Tx.Prepare(gamesSql)
	if err != nil {
		store.log.WithFields(logrus.Fields{
			"event":      "activeGamesStore::GetActiveGameDetails - Failed to prepare GetActiveGameDetails SELECT query.",
			"stackTrace": string(debug.Stack()),
		}).Error(err)
		return officialInfo, err
	}
	defer query.Close()

	officialGamesArray, err := getOfficialsGamesFromQuery(query, officialId)
	if err != nil {
		return officialInfo, err
	}

	var officialGamesSorted []models.OfficialGamesByLeague

	officialGamesMap := make(map[int]*models.OfficialGamesByLeague)

	// Iterate through the array and populate the map
	for _, game := range officialGamesArray {
		key := game.LeagueId
		if _, exists := officialGamesMap[key]; !exists {
			officialGamesMap[key] = &models.OfficialGamesByLeague{
				League:   game.LeagueName,
				LeagueId: game.LeagueId,
				Games:    []models.OfficialGameDetails{game},
			}
		} else {
			officialGamesMap[key].Games = append(officialGamesMap[key].Games, game)
		}
	}

	var leagues []string

	for _, value := range officialGamesMap {
		var mappedGamesByLeague models.OfficialGamesByLeague

		leagues = append(leagues, value.League)

		mappedGamesByLeague.League = value.League
		mappedGamesByLeague.LeagueId = value.LeagueId
		mappedGamesByLeague.Games = value.Games

		officialGamesSorted = append(officialGamesSorted, mappedGamesByLeague)
	}

	officialInfo.Games = officialGamesSorted

	statsSql := `SELECT
			    COUNT(DISTINCT gd.game_id) AS total_games,

			    ROUND(SUM(gd.overtime) / COUNT(DISTINCT gd.game_id) * 100, 2) AS overtime_percentage,
			    ROUND(SUM(gd.shootout) / COUNT(DISTINCT gd.game_id) * 100, 2) AS shootout_percentage,

			    AVG((CAST(SUBSTRING(gd.game_length, 1, POSITION(':' IN gd.game_length) - 1) AS SIGNED) * 60) +
      					CAST(SUBSTRING(gd.game_length, POSITION(':' IN gd.game_length) + 1) AS SIGNED)) AS avg_game_length_minutes,

			    AVG(gd.home_goals + gd.visitor_goals) AS avg_goals,
			    AVG(gd.home_goals) AS avg_home_goals,
			    AVG(gd.visitor_goals) AS avg_visitor_goals,

			    AVG(gd.home_powerplays) AS avg_home_powerplays,
			    AVG(gd.visitor_powerplays) AS avg_visitor_powerplays,

			    AVG(gd.home_pims) AS avg_home_pims,
			    AVG(gd.visitor_pims) AS avg_visitor_pims
		FROM officials_games og
		JOIN game_details gd ON og.game_id = gd.game_id
		WHERE og.official_id = ?;
	`

	var officialStats models.OfficialStats

	officialStats.Leagues = leagues

	err = store.database.Tx.QueryRow(statsSql, officialId).Scan(
		&officialStats.TotalGames,
		&officialStats.OvertimeAverage,
		&officialStats.ShootoutAverage,
		&officialStats.AverageGameTime,
		&officialStats.AverageGoalsPerGame,
		&officialStats.AverageHomeGoalsPerGame,
		&officialStats.AverageVisitorGoalsPerGame,
		&officialStats.AverageHomePowerPlaysPerGame,
		&officialStats.AverageVisitorPowerPlaysPerGame,
		&officialStats.AverageHomePimsPerGame,
		&officialStats.AverageVisitorPimsPerGame,
	)
	if err != nil {
		store.log.WithFields(logrus.Fields{
			"event":      "arenaStore::GetOfficialById - Failed to execute GetOfficialById SELECT query.",
			"stackTrace": string(debug.Stack()),
		}).Error(err)
		return officialInfo, err
	}

	input := "148.0000"
	minutes, err := strconv.ParseFloat(input, 64)
	if err != nil {
		store.log.WithFields(logrus.Fields{
			"event":      "arenaStore::GetOfficialById - Failed to parse float",
			"stackTrace": string(debug.Stack()),
		}).Error(err)
		return officialInfo, err
	}

	totalMinutes := int(math.Round(minutes))
	duration := time.Duration(totalMinutes) * time.Minute
	hours := int(duration.Hours())
	minutesRemaining := int(duration.Minutes()) % 60

	officialStats.AverageGameTime = fmt.Sprintf("%d:%02d", hours, minutesRemaining)

	officialInfo.OfficialStats = officialStats

	return officialInfo, nil
}

func getOfficialsGamesFromQuery(query *sql.Stmt, officialId int) ([]models.OfficialGameDetails, error) {
	rows, err := query.Query(officialId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	officialsGames := make([]models.OfficialGameDetails, 0)
	for rows.Next() {
		var officialGame models.OfficialGameDetails

		var refereeOne, refereeTwo, linesmanOne, linesmanTwo sql.NullString
		err = rows.Scan(
			&officialGame.GameId,
			&officialGame.LsGameId,
			&officialGame.LeagueName,
			&officialGame.LeagueId,
			&officialGame.HomeTeam,
			&officialGame.VisitorTeam,
			&officialGame.Arena,
			&officialGame.GameTime,
			&officialGame.GameDate,
			&refereeOne,
			&refereeTwo,
			&linesmanOne,
			&linesmanTwo,
		)
		if err != nil {
			return nil, err
		}

		if refereeOne.Valid {
			officialGame.RefereeOne = refereeOne.String
		} else {
			officialGame.RefereeOne = "Default Referee"
		}

		if refereeTwo.Valid {
			officialGame.RefereeTwo = refereeTwo.String
		} else {
			officialGame.RefereeTwo = "Default Referee"
		}

		if linesmanOne.Valid {
			officialGame.LinesmanOne = linesmanOne.String
		} else {
			officialGame.LinesmanOne = "Default Linesman"
		}

		if linesmanTwo.Valid {
			officialGame.LinesmanTwo = linesmanTwo.String
		} else {
			officialGame.LinesmanTwo = "Default Linesman"
		}

		officialsGames = append(officialsGames, officialGame)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return officialsGames, nil
}
