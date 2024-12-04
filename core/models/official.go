package models

type Official struct {
	Id           int
	FirstName    string
	LastName     string
	JerseyNumber int
}

type OfficialSearch struct {
	Id           int
	FirstName    string
	LastName     string
	JerseyNumber int
	Leagues      string
}

type OfficialInfo struct {
	Id            int
	FirstName     string
	LastName      string
	JerseyNumber  int
	Games         []OfficialGamesByLeague
	OfficialStats OfficialStats
}

type OfficialGamesByLeague struct {
	League   string
	LeagueId int
	Games    []OfficialGameDetails
}

type OfficialGameDetails struct {
	GameId      string
	LsGameId    string
	LeagueName  string
	LeagueId    int
	HomeTeam    string
	VisitorTeam string
	Arena       string
	GameDate    string
	GameTime    string
	RefereeOne  string
	RefereeTwo  string
	LinesmanOne string
	LinesmanTwo string
}

type OfficialStats struct {
	Leagues                         []string
	TotalGames                      int
	OvertimeAverage                 float32
	ShootoutAverage                 float32
	AverageGameTime                 string
	AverageGoalsPerGame             float32
	AverageHomeGoalsPerGame         float32
	AverageVisitorGoalsPerGame      float32
	AverageHomePowerPlaysPerGame    float32
	AverageVisitorPowerPlaysPerGame float32
	AverageHomePimsPerGame          float32
	AverageVisitorPimsPerGame       float32
}

type OfficialService interface {
	GetOfficialById(officialId int) (Official, error)
	GetOfficialInfoByOfficialId(officialId int) (OfficialInfo, error)
	GetOfficialsBySearchTerm(searchTerm string) ([]OfficialSearch, error)
}
