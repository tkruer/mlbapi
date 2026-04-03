package mlbapi

type TeamID int

type Team struct {
	ID            TeamID
	Name          string
	TeamCode      string
	FileCode      string
	Abbreviation  string
	TeamName      string
	LocationName  string
	ShortName     string
	FranchiseName string
	ClubName      string
}

type TeamReference struct {
	ID           TeamID
	Name         string
	Abbreviation string
}

type Position struct {
	Code         string
	Name         string
	Type         string
	Abbreviation string
}

type PlayerLookup struct {
	ID              int
	FullName        string
	FirstName       string
	LastName        string
	PrimaryNumber   string
	CurrentTeam     TeamReference
	PrimaryPosition Position
	Code            string
	Abbreviation    string
	UseName         string
	BoxscoreName    string
	Nickname        string
	MLBDebutDate    string
	NameFirstLast   string
	FirstLastName   string
	LastFirstName   string
	LastInitName    string
	InitLastName    string
	FullFMLName     string
	FullLFMName     string
	NameSlug        string
}

type ScheduleOptions struct {
	Date                string
	StartDate           string
	EndDate             string
	TeamID              TeamID
	OpponentID          TeamID
	SportID             int
	GameID              string
	LeagueID            string
	Season              int
	DisableSeriesStatus bool
}

type ScheduleGame struct {
	GameID              int
	GameDatetime        string
	GameDate            string
	GameType            string
	Status              string
	AwayName            string
	HomeName            string
	AwayID              TeamID
	HomeID              TeamID
	DoubleHeader        string
	GameNumber          int
	HomeProbablePitcher string
	AwayProbablePitcher string
	HomePitcherNote     string
	AwayPitcherNote     string
	AwayScore           int
	HomeScore           int
	CurrentInning       int
	InningState         string
	VenueID             int
	VenueName           string
	NationalBroadcasts  []string
	SeriesStatus        string
	WinningTeam         string
	LosingTeam          string
	WinningPitcher      string
	LosingPitcher       string
	SavePitcher         string
	Summary             string
}

type BoxscoreOptions struct {
	Timecode         string
	SkipBattingBox   bool
	SkipBattingInfo  bool
	SkipFieldingInfo bool
	SkipPitchingBox  bool
	SkipGameInfo     bool
}

type BatterLine struct {
	NameField    string
	AB           string
	R            string
	H            string
	Doubles      string
	Triples      string
	HR           string
	RBI          string
	SB           string
	BB           string
	K            string
	LOB          string
	AVG          string
	OPS          string
	PersonID     int
	Substitution bool
	Note         string
	Name         string
	Position     string
	OBP          string
	SLG          string
	BattingOrder string
}

type PitcherLine struct {
	NameField string
	IP        string
	H         string
	R         string
	ER        string
	BB        string
	K         string
	HR        string
	ERA       string
	P         string
	S         string
	Name      string
	PersonID  int
	Note      string
}

type InfoField struct {
	Label string
	Value string
}

type BoxscoreData struct {
	GameID             string
	TeamInfo           JSON
	PlayerInfo         JSON
	Away               JSON
	Home               JSON
	AwayBatters        []BatterLine
	HomeBatters        []BatterLine
	AwayBattingTotals  BatterLine
	HomeBattingTotals  BatterLine
	AwayBattingNotes   []string
	HomeBattingNotes   []string
	AwayPitchers       []PitcherLine
	HomePitchers       []PitcherLine
	AwayPitchingTotals PitcherLine
	HomePitchingTotals PitcherLine
	GameBoxInfo        []InfoField
}

type ScoringPlayData struct {
	Home  JSON
	Away  JSON
	Plays []JSON
}

type Playback struct {
	Name string
	URL  string
}

type HighlightItem struct {
	Date        string
	Title       string
	Headline    string
	Description string
	Duration    string
	Playbacks   []Playback
	Raw         JSON
}

type PlayerStatOptions struct {
	Group   string
	Type    string
	SportID int
	Season  int
}

type PlayerStatSplit struct {
	Type   string
	Group  string
	Season string
	Stats  JSON
}

type PlayerSummary struct {
	ID          int
	FirstName   string
	LastName    string
	Active      bool
	CurrentTeam string
	Position    string
	Nickname    string
	LastPlayed  string
	MLBDebut    string
	BatSide     string
	PitchHand   string
	Stats       []PlayerStatSplit
}

type LookupPlayerOptions struct {
	GameType string
	Season   int
	SportID  int
}

type LookupTeamOptions struct {
	ActiveStatus string
	Season       int
	SportIDs     string
}

type TeamLeaderOptions struct {
	Season          int
	LeaderGameTypes string
	Limit           int
}

type LeaderEntry struct {
	Rank  string
	Name  string
	Team  string
	Value string
}

type LeagueLeaderOptions struct {
	Season     int
	Limit      int
	StatGroup  string
	LeagueID   string
	GameTypes  string
	PlayerPool string
	SportID    int
	StatType   string
}

type StandingsOptions struct {
	LeagueID        string
	Division        string
	DisableWildcard bool
	Season          int
	StandingsTypes  string
	Date            string
}

type DivisionStandings struct {
	DivisionID   int
	DivisionName string
	Teams        []StandingTeam
}

type StandingTeam struct {
	Name                      string
	DivisionRank              string
	Wins                      int
	Losses                    int
	GamesBack                 string
	WildCardRank              string
	WildCardGamesBack         string
	WildCardEliminationNumber string
	EliminationNumber         string
	TeamID                    TeamID
	LeagueRank                string
	SportRank                 string
}

type RosterOptions struct {
	RosterType string
	Season     int
	Date       string
}
