package mlbapi

import (
	"fmt"
	"sort"
	"strings"
)

const (
	TeamLosAngelesAngels     TeamID = 108
	TeamArizonaDiamondbacks  TeamID = 109
	TeamBaltimoreOrioles     TeamID = 110
	TeamBostonRedSox         TeamID = 111
	TeamChicagoCubs          TeamID = 112
	TeamCincinnatiReds       TeamID = 113
	TeamClevelandGuardians   TeamID = 114
	TeamColoradoRockies      TeamID = 115
	TeamDetroitTigers        TeamID = 116
	TeamHoustonAstros        TeamID = 117
	TeamKansasCityRoyals     TeamID = 118
	TeamLosAngelesDodgers    TeamID = 119
	TeamWashingtonNationals  TeamID = 120
	TeamNewYorkMets          TeamID = 121
	TeamAthletics            TeamID = 133
	TeamPittsburghPirates    TeamID = 134
	TeamSanDiegoPadres       TeamID = 135
	TeamSeattleMariners      TeamID = 136
	TeamSanFranciscoGiants   TeamID = 137
	TeamStLouisCardinals     TeamID = 138
	TeamTampaBayRays         TeamID = 139
	TeamTexasRangers         TeamID = 140
	TeamTorontoBlueJays      TeamID = 141
	TeamMinnesotaTwins       TeamID = 142
	TeamPhiladelphiaPhillies TeamID = 143
	TeamAtlantaBraves        TeamID = 144
	TeamChicagoWhiteSox      TeamID = 145
	TeamMiamiMarlins         TeamID = 146
	TeamNewYorkYankees       TeamID = 147
	TeamMilwaukeeBrewers     TeamID = 158
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

func (id TeamID) Int() int {
	return int(id)
}

func (id TeamID) String() string {
	if team, ok := TeamByID(id); ok {
		return team.Name
	}
	return fmt.Sprintf("TeamID(%d)", int(id))
}

func AllTeams() []Team {
	out := make([]Team, len(knownTeams))
	copy(out, knownTeams)
	return out
}

func TeamByID(id TeamID) (Team, bool) {
	team, ok := teamByID[id]
	return team, ok
}

func TeamByAbbreviation(abbreviation string) (Team, bool) {
	team, ok := teamByAbbreviation[strings.ToUpper(strings.TrimSpace(abbreviation))]
	return team, ok
}

func TeamByName(name string) (Team, bool) {
	team, ok := teamByName[strings.ToLower(strings.TrimSpace(name))]
	return team, ok
}

func SortTeamsByID(teams []Team) {
	sort.Slice(teams, func(i, j int) bool {
		return teams[i].ID < teams[j].ID
	})
}
