package mlbapi

// TeamID is the stable MLB Stats API identifier for a team.
type TeamID int

// Team contains metadata about an MLB club.
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

// TeamReference is a lightweight team summary used in nested API responses.
type TeamReference struct {
	ID           TeamID
	Name         string
	Abbreviation string
}

// Position describes a baseball position in MLB Stats API responses.
type Position struct {
	Code         string
	Name         string
	Type         string
	Abbreviation string
}

// PlayerLookup is the typed result returned by LookupPlayer.
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

// ScheduleOptions controls Schedule queries.
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

// ScheduleGame is the typed schedule/game summary returned by Schedule.
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

// BoxscoreOptions controls which sections are rendered by Boxscore.
type BoxscoreOptions struct {
	Timecode         string
	SkipBattingBox   bool
	SkipBattingInfo  bool
	SkipFieldingInfo bool
	SkipPitchingBox  bool
	SkipGameInfo     bool
}

// BatterLine represents a single batter row in parsed boxscore output.
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

// PitcherLine represents a single pitcher row in parsed boxscore output.
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

// InfoField represents a label/value pair from game boxscore info.
type InfoField struct {
	Label string
	Value string
}

// BoxscoreData contains parsed boxscore data for a single game.
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

// ScoringPlayData contains home team data, away team data, and parsed scoring plays.
type ScoringPlayData struct {
	Home  JSON
	Away  JSON
	Plays []JSON
}

// Playback identifies a highlight playback rendition and URL.
type Playback struct {
	Name string
	URL  string
}

// HighlightItem is the typed result returned by GameHighlightData.
type HighlightItem struct {
	Date        string
	Title       string
	Headline    string
	Description string
	Duration    string
	Playbacks   []Playback
	Raw         JSON
}

// PlayerStatOptions controls PlayerStatData and PlayerStats queries.
type PlayerStatOptions struct {
	Group   string
	Type    string
	SportID int
	Season  int
}

// PlayerStatSplit represents one stat split in a player summary.
type PlayerStatSplit struct {
	Type   string
	Group  string
	Season string
	Stats  JSON
}

// PlayerSummary is the typed result returned by PlayerStatData.
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

// LookupPlayerOptions controls LookupPlayer queries.
type LookupPlayerOptions struct {
	GameType string
	Season   int
	SportID  int
}

// LookupTeamOptions controls LookupTeam queries.
type LookupTeamOptions struct {
	ActiveStatus string
	Season       int
	SportIDs     string
}

// TeamLeaderOptions controls TeamLeaderData and TeamLeaders queries.
type TeamLeaderOptions struct {
	Season          int
	LeaderGameTypes string
	Limit           int
}

// LeaderEntry represents a single leaderboard row.
type LeaderEntry struct {
	Rank  string
	Name  string
	Team  string
	Value string
}

// LeagueLeaderOptions controls LeagueLeaderData and LeagueLeaders queries.
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

// StandingsOptions controls StandingsData and Standings queries.
type StandingsOptions struct {
	LeagueID        string
	Division        string
	DisableWildcard bool
	Season          int
	StandingsTypes  string
	Date            string
}

// DivisionStandings groups team standings by division.
type DivisionStandings struct {
	DivisionID   int
	DivisionName string
	Teams        []StandingTeam
}

// StandingTeam represents one team row within division standings.
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

// RosterOptions controls roster queries and formatting.
type RosterOptions struct {
	RosterType string
	Season     int
	Date       string
}
