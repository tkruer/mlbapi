package mlbapi

import (
	"fmt"
	"sort"
	"strings"
)

// Team IDs that map to the current MLB clubs in the Stats API.
const (
	TeamLosAngelesAngels     TeamID = 108 // Los Angeles Angels.
	TeamArizonaDiamondbacks  TeamID = 109 // Arizona Diamondbacks.
	TeamBaltimoreOrioles     TeamID = 110 // Baltimore Orioles.
	TeamBostonRedSox         TeamID = 111 // Boston Red Sox.
	TeamChicagoCubs          TeamID = 112 // Chicago Cubs.
	TeamCincinnatiReds       TeamID = 113 // Cincinnati Reds.
	TeamClevelandGuardians   TeamID = 114 // Cleveland Guardians.
	TeamColoradoRockies      TeamID = 115 // Colorado Rockies.
	TeamDetroitTigers        TeamID = 116 // Detroit Tigers.
	TeamHoustonAstros        TeamID = 117 // Houston Astros.
	TeamKansasCityRoyals     TeamID = 118 // Kansas City Royals.
	TeamLosAngelesDodgers    TeamID = 119 // Los Angeles Dodgers.
	TeamWashingtonNationals  TeamID = 120 // Washington Nationals.
	TeamNewYorkMets          TeamID = 121 // New York Mets.
	TeamAthletics            TeamID = 133 // Athletics.
	TeamPittsburghPirates    TeamID = 134 // Pittsburgh Pirates.
	TeamSanDiegoPadres       TeamID = 135 // San Diego Padres.
	TeamSeattleMariners      TeamID = 136 // Seattle Mariners.
	TeamSanFranciscoGiants   TeamID = 137 // San Francisco Giants.
	TeamStLouisCardinals     TeamID = 138 // St. Louis Cardinals.
	TeamTampaBayRays         TeamID = 139 // Tampa Bay Rays.
	TeamTexasRangers         TeamID = 140 // Texas Rangers.
	TeamTorontoBlueJays      TeamID = 141 // Toronto Blue Jays.
	TeamMinnesotaTwins       TeamID = 142 // Minnesota Twins.
	TeamPhiladelphiaPhillies TeamID = 143 // Philadelphia Phillies.
	TeamAtlantaBraves        TeamID = 144 // Atlanta Braves.
	TeamChicagoWhiteSox      TeamID = 145 // Chicago White Sox.
	TeamMiamiMarlins         TeamID = 146 // Miami Marlins.
	TeamNewYorkYankees       TeamID = 147 // New York Yankees.
	TeamMilwaukeeBrewers     TeamID = 158 // Milwaukee Brewers.
)

var knownTeams = []Team{
	{ID: TeamLosAngelesAngels, Name: "Los Angeles Angels", TeamCode: "ana", FileCode: "ana", Abbreviation: "LAA", TeamName: "Angels", LocationName: "Anaheim", ShortName: "LA Angels", FranchiseName: "Los Angeles", ClubName: "Angels"},
	{ID: TeamArizonaDiamondbacks, Name: "Arizona Diamondbacks", TeamCode: "ari", FileCode: "ari", Abbreviation: "AZ", TeamName: "D-backs", LocationName: "Phoenix", ShortName: "Arizona", FranchiseName: "Arizona", ClubName: "Diamondbacks"},
	{ID: TeamBaltimoreOrioles, Name: "Baltimore Orioles", TeamCode: "bal", FileCode: "bal", Abbreviation: "BAL", TeamName: "Orioles", LocationName: "Baltimore", ShortName: "Baltimore", FranchiseName: "Baltimore", ClubName: "Orioles"},
	{ID: TeamBostonRedSox, Name: "Boston Red Sox", TeamCode: "bos", FileCode: "bos", Abbreviation: "BOS", TeamName: "Red Sox", LocationName: "Boston", ShortName: "Boston", FranchiseName: "Boston", ClubName: "Red Sox"},
	{ID: TeamChicagoCubs, Name: "Chicago Cubs", TeamCode: "chn", FileCode: "chc", Abbreviation: "CHC", TeamName: "Cubs", LocationName: "Chicago", ShortName: "Chi Cubs", FranchiseName: "Chicago", ClubName: "Cubs"},
	{ID: TeamCincinnatiReds, Name: "Cincinnati Reds", TeamCode: "cin", FileCode: "cin", Abbreviation: "CIN", TeamName: "Reds", LocationName: "Cincinnati", ShortName: "Cincinnati", FranchiseName: "Cincinnati", ClubName: "Reds"},
	{ID: TeamClevelandGuardians, Name: "Cleveland Guardians", TeamCode: "cle", FileCode: "cle", Abbreviation: "CLE", TeamName: "Guardians", LocationName: "Cleveland", ShortName: "Cleveland", FranchiseName: "Cleveland", ClubName: "Guardians"},
	{ID: TeamColoradoRockies, Name: "Colorado Rockies", TeamCode: "col", FileCode: "col", Abbreviation: "COL", TeamName: "Rockies", LocationName: "Denver", ShortName: "Colorado", FranchiseName: "Colorado", ClubName: "Rockies"},
	{ID: TeamDetroitTigers, Name: "Detroit Tigers", TeamCode: "det", FileCode: "det", Abbreviation: "DET", TeamName: "Tigers", LocationName: "Detroit", ShortName: "Detroit", FranchiseName: "Detroit", ClubName: "Tigers"},
	{ID: TeamHoustonAstros, Name: "Houston Astros", TeamCode: "hou", FileCode: "hou", Abbreviation: "HOU", TeamName: "Astros", LocationName: "Houston", ShortName: "Houston", FranchiseName: "Houston", ClubName: "Astros"},
	{ID: TeamKansasCityRoyals, Name: "Kansas City Royals", TeamCode: "kca", FileCode: "kc", Abbreviation: "KC", TeamName: "Royals", LocationName: "Kansas City", ShortName: "Kansas City", FranchiseName: "Kansas City", ClubName: "Royals"},
	{ID: TeamLosAngelesDodgers, Name: "Los Angeles Dodgers", TeamCode: "lan", FileCode: "la", Abbreviation: "LAD", TeamName: "Dodgers", LocationName: "Los Angeles", ShortName: "LA Dodgers", FranchiseName: "Los Angeles", ClubName: "Dodgers"},
	{ID: TeamWashingtonNationals, Name: "Washington Nationals", TeamCode: "was", FileCode: "was", Abbreviation: "WSH", TeamName: "Nationals", LocationName: "Washington", ShortName: "Washington", FranchiseName: "Washington", ClubName: "Nationals"},
	{ID: TeamNewYorkMets, Name: "New York Mets", TeamCode: "nyn", FileCode: "nym", Abbreviation: "NYM", TeamName: "Mets", LocationName: "Flushing", ShortName: "NY Mets", FranchiseName: "New York", ClubName: "Mets"},
	{ID: TeamAthletics, Name: "Athletics", TeamCode: "ath", FileCode: "ath", Abbreviation: "ATH", TeamName: "Athletics", LocationName: "Sacramento", ShortName: "Athletics", FranchiseName: "Athletics", ClubName: "Athletics"},
	{ID: TeamPittsburghPirates, Name: "Pittsburgh Pirates", TeamCode: "pit", FileCode: "pit", Abbreviation: "PIT", TeamName: "Pirates", LocationName: "Pittsburgh", ShortName: "Pittsburgh", FranchiseName: "Pittsburgh", ClubName: "Pirates"},
	{ID: TeamSanDiegoPadres, Name: "San Diego Padres", TeamCode: "sdn", FileCode: "sd", Abbreviation: "SD", TeamName: "Padres", LocationName: "San Diego", ShortName: "San Diego", FranchiseName: "San Diego", ClubName: "Padres"},
	{ID: TeamSeattleMariners, Name: "Seattle Mariners", TeamCode: "sea", FileCode: "sea", Abbreviation: "SEA", TeamName: "Mariners", LocationName: "Seattle", ShortName: "Seattle", FranchiseName: "Seattle", ClubName: "Mariners"},
	{ID: TeamSanFranciscoGiants, Name: "San Francisco Giants", TeamCode: "sfn", FileCode: "sf", Abbreviation: "SF", TeamName: "Giants", LocationName: "San Francisco", ShortName: "San Francisco", FranchiseName: "San Francisco", ClubName: "Giants"},
	{ID: TeamStLouisCardinals, Name: "St. Louis Cardinals", TeamCode: "sln", FileCode: "stl", Abbreviation: "STL", TeamName: "Cardinals", LocationName: "St. Louis", ShortName: "St. Louis", FranchiseName: "St. Louis", ClubName: "Cardinals"},
	{ID: TeamTampaBayRays, Name: "Tampa Bay Rays", TeamCode: "tba", FileCode: "tb", Abbreviation: "TB", TeamName: "Rays", LocationName: "St. Petersburg", ShortName: "Tampa Bay", FranchiseName: "Tampa Bay", ClubName: "Rays"},
	{ID: TeamTexasRangers, Name: "Texas Rangers", TeamCode: "tex", FileCode: "tex", Abbreviation: "TEX", TeamName: "Rangers", LocationName: "Arlington", ShortName: "Texas", FranchiseName: "Texas", ClubName: "Rangers"},
	{ID: TeamTorontoBlueJays, Name: "Toronto Blue Jays", TeamCode: "tor", FileCode: "tor", Abbreviation: "TOR", TeamName: "Blue Jays", LocationName: "Toronto", ShortName: "Toronto", FranchiseName: "Toronto", ClubName: "Blue Jays"},
	{ID: TeamMinnesotaTwins, Name: "Minnesota Twins", TeamCode: "min", FileCode: "min", Abbreviation: "MIN", TeamName: "Twins", LocationName: "Minneapolis", ShortName: "Minnesota", FranchiseName: "Minnesota", ClubName: "Twins"},
	{ID: TeamPhiladelphiaPhillies, Name: "Philadelphia Phillies", TeamCode: "phi", FileCode: "phi", Abbreviation: "PHI", TeamName: "Phillies", LocationName: "Philadelphia", ShortName: "Philadelphia", FranchiseName: "Philadelphia", ClubName: "Phillies"},
	{ID: TeamAtlantaBraves, Name: "Atlanta Braves", TeamCode: "atl", FileCode: "atl", Abbreviation: "ATL", TeamName: "Braves", LocationName: "Atlanta", ShortName: "Atlanta", FranchiseName: "Atlanta", ClubName: "Braves"},
	{ID: TeamChicagoWhiteSox, Name: "Chicago White Sox", TeamCode: "cha", FileCode: "cws", Abbreviation: "CWS", TeamName: "White Sox", LocationName: "Chicago", ShortName: "Chi White Sox", FranchiseName: "Chicago", ClubName: "White Sox"},
	{ID: TeamMiamiMarlins, Name: "Miami Marlins", TeamCode: "mia", FileCode: "mia", Abbreviation: "MIA", TeamName: "Marlins", LocationName: "Miami", ShortName: "Miami", FranchiseName: "Miami", ClubName: "Marlins"},
	{ID: TeamNewYorkYankees, Name: "New York Yankees", TeamCode: "nya", FileCode: "nyy", Abbreviation: "NYY", TeamName: "Yankees", LocationName: "Bronx", ShortName: "NY Yankees", FranchiseName: "New York", ClubName: "Yankees"},
	{ID: TeamMilwaukeeBrewers, Name: "Milwaukee Brewers", TeamCode: "mil", FileCode: "mil", Abbreviation: "MIL", TeamName: "Brewers", LocationName: "Milwaukee", ShortName: "Milwaukee", FranchiseName: "Milwaukee", ClubName: "Brewers"},
}

var teamByID = func() map[TeamID]Team {
	index := make(map[TeamID]Team, len(knownTeams))
	for _, team := range knownTeams {
		index[team.ID] = team
	}
	return index
}()

var teamByAbbreviation = func() map[string]Team {
	index := make(map[string]Team, len(knownTeams))
	for _, team := range knownTeams {
		index[strings.ToUpper(team.Abbreviation)] = team
	}
	return index
}()

var teamByName = func() map[string]Team {
	index := make(map[string]Team, len(knownTeams)*4)
	for _, team := range knownTeams {
		for _, key := range []string{
			team.Name,
			team.TeamName,
			team.ShortName,
			strings.TrimSpace(team.FranchiseName + " " + team.ClubName),
		} {
			if key != "" {
				index[strings.ToLower(key)] = team
			}
		}
	}
	return index
}()

// Int returns the raw integer value used by the MLB Stats API.
func (id TeamID) Int() int {
	return int(id)
}

// String returns the full team name for known IDs.
func (id TeamID) String() string {
	if team, ok := TeamByID(id); ok {
		return team.Name
	}
	return fmt.Sprintf("TeamID(%d)", int(id))
}

// AllTeams returns a copy of the built-in team metadata list.
func AllTeams() []Team {
	out := make([]Team, len(knownTeams))
	copy(out, knownTeams)
	return out
}

// TeamByID looks up a team by its stable TeamID.
func TeamByID(id TeamID) (Team, bool) {
	team, ok := teamByID[id]
	return team, ok
}

// TeamByAbbreviation looks up a team by its MLB abbreviation.
func TeamByAbbreviation(abbreviation string) (Team, bool) {
	team, ok := teamByAbbreviation[strings.ToUpper(strings.TrimSpace(abbreviation))]
	return team, ok
}

// TeamByName looks up a team by common display names and aliases.
func TeamByName(name string) (Team, bool) {
	team, ok := teamByName[strings.ToLower(strings.TrimSpace(name))]
	return team, ok
}

// SortTeamsByID sorts teams in place by ascending TeamID.
func SortTeamsByID(teams []Team) {
	sort.Slice(teams, func(i, j int) bool {
		return teams[i].ID < teams[j].ID
	})
}
