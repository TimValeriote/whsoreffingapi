# whosreffing-api

## Whosreffing API docs
### Base URL for development:
```
http://localhost:3095/whosreffing
```
## Routes:

### /activegames
-- Gets all active games currently available

**type:** GET

**response:**

```json
{
    "active_games": [
        {
            "id": int,
            "game_id": int,
            "ls_game_id": int
        },
        {
            "id": int,
            "game_id": int,
            "ls_game_id": int
        }
    ]
}
```

### /arena/:arenaId
-- Gets a specific arena by arena id

**type:** GET

**url params:**

`arenaId: int`

**response:**

```json
{
    "arena": {
        "id": int,
        "name": string
    }
}
```

### /game/:gameId
-- Gets a single game by a game id

**type:** GET

**url params:**

`gameId: int`

**response:**

```json
{
    "game": {
        "id": int,
        "ls_game_id": int,
        "league_id": int,
        "date": string(date) ex: "2023-10-11T00:00:00Z",
        "time": string(time) ex: "19:45:00",
        "home_team_id": int,
        "visiting_team_id": int,
        "arena_id": int,
        "status_id": int,
        "referee_one_id": int,
        "referee_two_id": int,
        "linesman_one_id": int,
        "linesman_two_id": int,
        "season_id": int
    }
}
```

### /gamedetails/:gameId
-- Gets a games details by a game id

**type:** GET

**url params:**

`gameId: int`

**response:**

```json
{
    "game_details": {
        "id": int,
        "game_id": int,
        "game_length": string ex: "2:30",
        "home_goals": int,
        "visiting_goals": int,
        "overtime": int (1 or 0),
        "shootout": int (1 or 0),
        "home_power_plays": int,
        "visitor_power_plays": int,
        "home_pims": int,
        "visitor_pims": int,
        "home_faceoffs": int,
        "visitor_faceoffs": int,
        "home_faceoffs_won": int,
        "visitor_faceoffs_won": int
    }
}
```

### /official/:officialId
-- Gets a single official by their id

**type:** GET

**url params:**

`officialId: int`

**response:**

```json
{
    "official": {
        "id": int,
        "first_name": string,
        "last_name": string,
        "jersey_number": int
    }
}
```

### /officialgames/:officialId
-- Gets all games for an official

**type:** GET

**url params:**

`officialId: int`

**response:**

```json
{
    "official_games": [
        {
            "id": int,
            "game_id": int,
            "official_id": int
        },
        {
            "id": int,
            "game_id": int,
            "official_id": int
        }
    ]
}
```

### /leagues
-- Gets all leagues

**type:** GET

**response:**

```json
{
    "leagues": [
        {
            "id": int,
            "name": string,
            "level_id": int,
            "code": string
        },
        {
            "id": int,
            "name": string,
            "level_id": int,
            "code": string
        }
    ]
}
```


### /league/:leagueId
-- Gets a single league by id

**type:** GET

**url params:**

`leagueId: int`

**response:**

```json
{
    "league": {
        "id": int,
        "name": string,
        "level_id": int,
        "code": string
    }
}
```
