package mlbapi

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"
)

var validMetaTypes = map[string]struct{}{
	"awards":             {},
	"baseballStats":      {},
	"eventTypes":         {},
	"freeGameTypes":      {},
	"gameStatus":         {},
	"gameTypes":          {},
	"hitTrajectories":    {},
	"jobTypes":           {},
	"languages":          {},
	"leagueLeaderTypes":  {},
	"logicalEvents":      {},
	"metrics":            {},
	"pitchCodes":         {},
	"pitchTypes":         {},
	"platforms":          {},
	"positions":          {},
	"reviewReasons":      {},
	"rosterTypes":        {},
	"runnerDetailTypes":  {},
	"scheduleTypes":      {},
	"scheduleEventTypes": {},
	"situationCodes":     {},
	"sky":                {},
	"standingsTypes":     {},
	"statGroups":         {},
	"statTypes":          {},
	"violationTypes":     {},
	"windDirection":      {},
}

func GamePace(ctx context.Context, season, sportID int) (string, error) {
	return DefaultClient.GamePace(ctx, season, sportID)
}

func (c *Client) GamePace(ctx context.Context, season, sportID int) (string, error) {
	data, err := c.GamePaceData(ctx, season, sportID)
	if err != nil {
		return "", err
	}

	if season == 0 {
		season = currentYear()
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%d Game Pace Stats\n", season))
	for _, sportValue := range nestedSlice(data, "sports") {
		sport := asMap(sportValue)
		for _, key := range sortedKeys(sport) {
			if key == "season" || key == "sport" {
				continue
			}
			if key == "prPortalCalculatedFields" {
				fields := asMap(sport[key])
				for _, nestedKey := range sortedKeys(fields) {
					builder.WriteString(fmt.Sprintf("%s: %s\n", nestedKey, stringValue(fields[nestedKey])))
				}
				continue
			}
			builder.WriteString(fmt.Sprintf("%s: %s\n", key, stringValue(sport[key])))
		}
	}

	return builder.String(), nil
}

func GamePaceData(ctx context.Context, season, sportID int) (JSON, error) {
	return DefaultClient.GamePaceData(ctx, season, sportID)
}

func (c *Client) GamePaceData(ctx context.Context, season, sportID int) (JSON, error) {
	params := Params{}
	if season == 0 {
		season = currentYear()
	}
	if sportID == 0 {
		sportID = 1
	}

	params["season"] = season
	params["sportId"] = sportID

	response, err := c.CallGamePace(ctx, params)
	if err != nil {
		return nil, err
	}
	if len(nestedSlice(response, "sports")) == 0 {
		return nil, fmt.Errorf("mlbapi: no game pace info found for season %d; data appears to begin in 1999", season)
	}
	return response, nil
}

func PlayerStats(ctx context.Context, personID int, opts PlayerStatOptions) (string, error) {
	return DefaultClient.PlayerStats(ctx, personID, opts)
}

func (c *Client) PlayerStats(ctx context.Context, personID int, opts PlayerStatOptions) (string, error) {
	player, err := c.PlayerStatData(ctx, personID, opts)
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	builder.WriteString(player.FirstName)
	if player.Nickname != "" {
		builder.WriteString(fmt.Sprintf(" %q", player.Nickname))
	}
	builder.WriteString(fmt.Sprintf(" %s, %s (%s-", player.LastName, player.Position, yearPrefix(player.MLBDebut)))
	if !player.Active {
		builder.WriteString(yearPrefix(player.LastPlayed))
	}
	builder.WriteString(")\n\n")

	for _, statGroup := range player.Stats {
		builder.WriteString(statGroup.Type + " " + statGroup.Group)
		position := nestedString(statGroup.Stats, "position", "abbreviation")
		if position != "" {
			builder.WriteString(fmt.Sprintf(" (%s)", position))
		}
		builder.WriteByte('\n')
		for _, key := range sortedKeys(statGroup.Stats) {
			if key == "position" {
				continue
			}
			builder.WriteString(fmt.Sprintf("%s: %s\n", key, stringValue(statGroup.Stats[key])))
		}
		builder.WriteByte('\n')
	}

	return builder.String(), nil
}

func PlayerStatData(ctx context.Context, personID int, opts PlayerStatOptions) (PlayerSummary, error) {
	return DefaultClient.PlayerStatData(ctx, personID, opts)
}

func (c *Client) PlayerStatData(ctx context.Context, personID int, opts PlayerStatOptions) (PlayerSummary, error) {
	if opts.Type == "" {
		opts.Type = "season"
	}
	if opts.Group == "" {
		opts.Group = "[hitting,pitching,fielding]"
	}
	if opts.SportID == 0 {
		opts.SportID = 1
	}
	if opts.Season != 0 && !strings.Contains(opts.Type, "season") {
		return PlayerSummary{}, fmt.Errorf("mlbapi: season is only valid when using a season stat type")
	}

	hydrate := fmt.Sprintf("stats(group=%s,type=%s", opts.Group, opts.Type)
	if opts.Season != 0 {
		hydrate += fmt.Sprintf(",season=%d", opts.Season)
	}
	hydrate += fmt.Sprintf(",sportId=%d),currentTeam", opts.SportID)

	response, err := c.CallPerson(ctx, Params{
		"personId": personID,
		"hydrate":  hydrate,
	})
	if err != nil {
		return PlayerSummary{}, err
	}

	people := nestedSlice(response, "people")
	if len(people) == 0 {
		return PlayerSummary{}, fmt.Errorf("mlbapi: player %d not found", personID)
	}

	person := asMap(people[0])
	player := PlayerSummary{
		ID:          intValue(person["id"]),
		FirstName:   stringValue(person["useName"]),
		LastName:    stringValue(person["lastName"]),
		Active:      boolValue(person["active"]),
		CurrentTeam: nestedString(person, "currentTeam", "name"),
		Position:    nestedString(person, "primaryPosition", "abbreviation"),
		Nickname:    stringValue(person["nickName"]),
		LastPlayed:  stringValue(person["lastPlayedDate"]),
		MLBDebut:    stringValue(person["mlbDebutDate"]),
		BatSide:     nestedString(person, "batSide", "description"),
		PitchHand:   nestedString(person, "pitchHand", "description"),
	}

	for _, statValue := range asSlice(person["stats"]) {
		stat := asMap(statValue)
		for _, splitValue := range asSlice(stat["splits"]) {
			split := asMap(splitValue)
			player.Stats = append(player.Stats, PlayerStatSplit{
				Type:   nestedString(stat, "type", "displayName"),
				Group:  nestedString(stat, "group", "displayName"),
				Season: stringValue(split["season"]),
				Stats:  asMap(split["stat"]),
			})
		}
	}

	return player, nil
}

func LatestSeason(ctx context.Context, sportID int) (JSON, error) {
	return DefaultClient.LatestSeason(ctx, sportID)
}

func (c *Client) LatestSeason(ctx context.Context, sportID int) (JSON, error) {
	if sportID == 0 {
		sportID = 1
	}

	response, err := c.CallSeason(ctx, Params{
		"sportId":  sportID,
		"seasonId": "all",
	})
	if err != nil {
		return nil, err
	}

	seasons := nestedSlice(response, "seasons")
	if len(seasons) == 0 {
		return nil, fmt.Errorf("mlbapi: no seasons found for sport %d", sportID)
	}

	today := time.Now().Format("2006-01-02")
	for _, seasonValue := range seasons {
		season := asMap(seasonValue)
		if today < stringValue(season["seasonEndDate"]) {
			return season, nil
		}
	}

	return asMap(seasons[len(seasons)-1]), nil
}

func LookupPlayer(ctx context.Context, lookup string, opts LookupPlayerOptions) ([]PlayerLookup, error) {
	return DefaultClient.LookupPlayer(ctx, lookup, opts)
}

func (c *Client) LookupPlayer(ctx context.Context, lookup string, opts LookupPlayerOptions) ([]PlayerLookup, error) {
	if opts.SportID == 0 {
		opts.SportID = 1
	}

	params := Params{
		"sportId": opts.SportID,
		"fields":  "people,id,fullName,firstName,lastName,primaryNumber,currentTeam,id,name,abbreviation,primaryPosition,code,name,type,abbreviation,useName,boxscoreName,nickName,mlbDebutDate,nameFirstLast,firstLastName,lastFirstName,lastInitName,initLastName,fullFMLName,fullLFMName,nameSlug",
	}
	if opts.GameType != "" {
		params["gameType"] = opts.GameType
	}
	if opts.Season == 0 {
		season, err := c.LatestSeason(ctx, opts.SportID)
		if err != nil {
			return nil, err
		}
		opts.Season = intValue(season["seasonId"])
	}
	params["season"] = opts.Season

	response, err := c.CallSportsPlayers(ctx, params)
	if err != nil {
		return nil, err
	}

	terms := strings.Fields(strings.ToLower(lookup))
	players := make([]PlayerLookup, 0)
	for _, playerValue := range nestedSlice(response, "people") {
		player := playerLookupFromJSON(asMap(playerValue))
		if playerLookupMatches(player, terms) {
			players = append(players, player)
		}
	}

	return players, nil
}

func LookupTeam(ctx context.Context, lookup string, opts LookupTeamOptions) ([]Team, error) {
	return DefaultClient.LookupTeam(ctx, lookup, opts)
}

func (c *Client) LookupTeam(ctx context.Context, lookup string, opts LookupTeamOptions) ([]Team, error) {
	if opts.ActiveStatus == "" {
		opts.ActiveStatus = "Y"
	}
	if opts.SportIDs == "" {
		opts.SportIDs = "1"
	}

	params := Params{
		"activeStatus": opts.ActiveStatus,
		"sportIds":     opts.SportIDs,
		"fields":       "teams,id,name,teamCode,fileCode,abbreviation,teamName,locationName,shortName,franchiseName,clubName",
	}
	if opts.Season == 0 {
		season, err := c.LatestSeason(ctx, intValue(strings.Split(opts.SportIDs, ",")[0]))
		if err != nil {
			return nil, err
		}
		opts.Season = intValue(season["seasonId"])
	}
	params["season"] = opts.Season

	response, err := c.CallTeams(ctx, params)
	if err != nil {
		return nil, err
	}

	target := strings.ToLower(lookup)
	teams := make([]Team, 0)
	for _, teamValue := range nestedSlice(response, "teams") {
		team := teamFromJSON(asMap(teamValue))
		if teamMatches(team, target) {
			teams = append(teams, team)
		}
	}

	return teams, nil
}

func TeamLeaders(ctx context.Context, teamID TeamID, leaderCategories string, opts TeamLeaderOptions) (string, error) {
	return DefaultClient.TeamLeaders(ctx, teamID, leaderCategories, opts)
}

func (c *Client) TeamLeaders(ctx context.Context, teamID TeamID, leaderCategories string, opts TeamLeaderOptions) (string, error) {
	lines, err := c.TeamLeaderData(ctx, teamID, leaderCategories, opts)
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%-4s %-20s %-5s\n", "Rank", "Name", "Value"))
	for _, line := range lines {
		builder.WriteString(fmt.Sprintf("%4s %-20s %5s\n", line.Rank, line.Name, line.Value))
	}
	return builder.String(), nil
}

func TeamLeaderData(ctx context.Context, teamID TeamID, leaderCategories string, opts TeamLeaderOptions) ([]LeaderEntry, error) {
	return DefaultClient.TeamLeaderData(ctx, teamID, leaderCategories, opts)
}

func (c *Client) TeamLeaderData(ctx context.Context, teamID TeamID, leaderCategories string, opts TeamLeaderOptions) ([]LeaderEntry, error) {
	if opts.Season == 0 {
		opts.Season = currentYear()
	}
	if opts.LeaderGameTypes == "" {
		opts.LeaderGameTypes = "R"
	}
	if opts.Limit == 0 {
		opts.Limit = 10
	}

	response, err := c.CallTeamLeaders(ctx, Params{
		"leaderCategories": leaderCategories,
		"season":           opts.Season,
		"teamId":           teamID,
		"leaderGameTypes":  opts.LeaderGameTypes,
		"limit":            opts.Limit,
		"fields":           "teamLeaders,leaders,rank,value,person,fullName",
	})
	if err != nil {
		return nil, err
	}

	leaders := make([]LeaderEntry, 0)
	sections := nestedSlice(response, "teamLeaders")
	if len(sections) == 0 {
		return leaders, nil
	}
	for _, leaderValue := range nestedSlice(asMap(sections[0]), "leaders") {
		leader := asMap(leaderValue)
		leaders = append(leaders, LeaderEntry{
			Rank:  stringValue(leader["rank"]),
			Name:  nestedString(leader, "person", "fullName"),
			Value: stringValue(leader["value"]),
		})
	}

	return leaders, nil
}

func LeagueLeaders(ctx context.Context, leaderCategories string, opts LeagueLeaderOptions) (string, error) {
	return DefaultClient.LeagueLeaders(ctx, leaderCategories, opts)
}

func (c *Client) LeagueLeaders(ctx context.Context, leaderCategories string, opts LeagueLeaderOptions) (string, error) {
	lines, err := c.LeagueLeaderData(ctx, leaderCategories, opts)
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%-4s %-20s %-23s %-5s\n", "Rank", "Name", "Team", "Value"))
	for _, line := range lines {
		builder.WriteString(fmt.Sprintf("%4s %-20s %-23s %5s\n", line.Rank, line.Name, line.Team, line.Value))
	}
	return builder.String(), nil
}

func LeagueLeaderData(ctx context.Context, leaderCategories string, opts LeagueLeaderOptions) ([]LeaderEntry, error) {
	return DefaultClient.LeagueLeaderData(ctx, leaderCategories, opts)
}

func (c *Client) LeagueLeaderData(ctx context.Context, leaderCategories string, opts LeagueLeaderOptions) ([]LeaderEntry, error) {
	if opts.Limit == 0 {
		opts.Limit = 10
	}
	if opts.SportID == 0 {
		opts.SportID = 1
	}

	params := Params{
		"leaderCategories": leaderCategories,
		"sportId":          opts.SportID,
		"limit":            opts.Limit,
	}
	if opts.Season != 0 {
		params["season"] = opts.Season
	}
	if opts.StatType != "" {
		params["statType"] = opts.StatType
	}
	if opts.Season == 0 && opts.StatType == "" {
		params["season"] = currentYear()
	}
	if opts.StatGroup != "" {
		if opts.StatGroup == "batting" {
			opts.StatGroup = "hitting"
		}
		params["statGroup"] = opts.StatGroup
	}
	if opts.GameTypes != "" {
		params["leaderGameTypes"] = opts.GameTypes
	}
	if opts.LeagueID != "" {
		params["leagueId"] = opts.LeagueID
	}
	if opts.PlayerPool != "" {
		params["playerPool"] = opts.PlayerPool
	}
	params["fields"] = "leagueLeaders,leaders,rank,value,team,name,league,name,person,fullName"

	response, err := c.CallStatsLeaders(ctx, params)
	if err != nil {
		return nil, err
	}

	leaders := make([]LeaderEntry, 0)
	sections := nestedSlice(response, "leagueLeaders")
	if len(sections) == 0 {
		return leaders, nil
	}
	for _, leaderValue := range nestedSlice(asMap(sections[0]), "leaders") {
		leader := asMap(leaderValue)
		leaders = append(leaders, LeaderEntry{
			Rank:  stringValue(leader["rank"]),
			Name:  nestedString(leader, "person", "fullName"),
			Team:  nestedString(leader, "team", "name"),
			Value: stringValue(leader["value"]),
		})
	}

	return leaders, nil
}

func Standings(ctx context.Context, opts StandingsOptions) (string, error) {
	return DefaultClient.Standings(ctx, opts)
}

func (c *Client) Standings(ctx context.Context, opts StandingsOptions) (string, error) {
	divisions, err := c.StandingsData(ctx, opts)
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	for _, division := range divisions {
		builder.WriteString(division.DivisionName)
		builder.WriteByte('\n')
		if !opts.DisableWildcard {
			builder.WriteString(fmt.Sprintf("%4s %-21s %3s %3s %4s %4s %7s %5s %4s\n", "Rank", "Team", "W", "L", "GB", "(E#)", "WC Rank", "WC GB", "(E#)"))
			for _, team := range division.Teams {
				builder.WriteString(fmt.Sprintf("%4s %-21s %3d %3d %4s %4s %7s %5s %4s\n",
					team.DivisionRank,
					team.Name,
					team.Wins,
					team.Losses,
					team.GamesBack,
					team.EliminationNumber,
					team.WildCardRank,
					team.WildCardGamesBack,
					team.WildCardEliminationNumber,
				))
			}
		} else {
			builder.WriteString(fmt.Sprintf("%4s %-21s %3s %3s %4s %4s\n", "Rank", "Team", "W", "L", "GB", "(E#)"))
			for _, team := range division.Teams {
				builder.WriteString(fmt.Sprintf("%4s %-21s %3d %3d %4s %4s\n",
					team.DivisionRank,
					team.Name,
					team.Wins,
					team.Losses,
					team.GamesBack,
					team.EliminationNumber,
				))
			}
		}
		builder.WriteByte('\n')
	}

	return builder.String(), nil
}

func StandingsData(ctx context.Context, opts StandingsOptions) ([]DivisionStandings, error) {
	return DefaultClient.StandingsData(ctx, opts)
}

func (c *Client) StandingsData(ctx context.Context, opts StandingsOptions) ([]DivisionStandings, error) {
	if opts.LeagueID == "" {
		opts.LeagueID = "103,104"
	}
	if opts.Division == "" {
		opts.Division = "all"
	}
	if opts.Date != "" {
		opts.Season = seasonFromDate(opts.Date)
	}
	if opts.Season == 0 {
		opts.Season = currentYear()
	}
	if opts.StandingsTypes == "" {
		opts.StandingsTypes = "regularSeason"
	}

	params := Params{
		"leagueId":       opts.LeagueID,
		"season":         opts.Season,
		"standingsTypes": opts.StandingsTypes,
		"hydrate":        "team(division)",
		"fields":         "records,standingsType,teamRecords,team,name,division,id,nameShort,abbreviation,divisionRank,gamesBack,wildCardRank,wildCardGamesBack,wildCardEliminationNumber,divisionGamesBack,clinched,eliminationNumber,winningPercentage,type,wins,losses,leagueRank,sportRank",
	}
	if opts.Date != "" {
		params["date"] = opts.Date
	}

	response, err := c.CallStandings(ctx, params)
	if err != nil {
		return nil, err
	}

	divisions := make([]DivisionStandings, 0)
	indexByID := make(map[int]int)
	divisionFilter := strings.ToLower(opts.Division)

	for _, recordValue := range nestedSlice(response, "records") {
		record := asMap(recordValue)
		for _, teamValue := range asSlice(record["teamRecords"]) {
			teamRecord := asMap(teamValue)
			division := nestedMap(teamRecord, "team", "division")
			divisionID := intValue(division["id"])
			divisionAbbrev := strings.ToLower(stringValue(division["abbreviation"]))
			if divisionFilter != "all" && divisionFilter != divisionAbbrev && divisionFilter != stringValue(division["id"]) {
				continue
			}

			index, ok := indexByID[divisionID]
			if !ok {
				index = len(divisions)
				indexByID[divisionID] = index
				divisions = append(divisions, DivisionStandings{
					DivisionID:   divisionID,
					DivisionName: stringValue(division["name"]),
				})
			}

			divisions[index].Teams = append(divisions[index].Teams, StandingTeam{
				Name:                      nestedString(teamRecord, "team", "name"),
				DivisionRank:              stringValue(teamRecord["divisionRank"]),
				Wins:                      intValue(teamRecord["wins"]),
				Losses:                    intValue(teamRecord["losses"]),
				GamesBack:                 stringValue(teamRecord["gamesBack"]),
				WildCardRank:              defaultString(teamRecord, "wildCardRank", "-"),
				WildCardGamesBack:         defaultString(teamRecord, "wildCardGamesBack", "-"),
				WildCardEliminationNumber: defaultString(teamRecord, "wildCardEliminationNumber", "-"),
				EliminationNumber:         defaultString(teamRecord, "eliminationNumber", "-"),
				TeamID:                    TeamID(intValue(nestedValue(teamRecord, "team", "id"))),
				LeagueRank:                defaultString(teamRecord, "leagueRank", "-"),
				SportRank:                 defaultString(teamRecord, "sportRank", "-"),
			})
		}
	}

	return divisions, nil
}

func Roster(ctx context.Context, teamID TeamID, opts RosterOptions) (string, error) {
	return DefaultClient.Roster(ctx, teamID, opts)
}

func (c *Client) Roster(ctx context.Context, teamID TeamID, opts RosterOptions) (string, error) {
	if opts.RosterType == "" {
		opts.RosterType = "active"
	}
	if opts.Season == 0 {
		opts.Season = currentYear()
	}

	params := Params{
		"rosterType": opts.RosterType,
		"season":     opts.Season,
		"teamId":     teamID,
	}
	if opts.Date != "" {
		params["date"] = opts.Date
	}

	response, err := c.CallTeamRoster(ctx, params)
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	for _, playerValue := range nestedSlice(response, "roster") {
		player := asMap(playerValue)
		builder.WriteString(fmt.Sprintf("#%-3s %-3s %s\n",
			stringValue(player["jerseyNumber"]),
			nestedString(player, "position", "abbreviation"),
			nestedString(player, "person", "fullName"),
		))
	}

	return builder.String(), nil
}

func Meta(ctx context.Context, metaType string) (JSON, error) {
	return DefaultClient.Meta(ctx, metaType)
}

func (c *Client) Meta(ctx context.Context, metaType string) (JSON, error) {
	if _, ok := validMetaTypes[metaType]; !ok {
		valid := make([]string, 0, len(validMetaTypes))
		for key := range validMetaTypes {
			valid = append(valid, key)
		}
		sort.Strings(valid)
		return nil, fmt.Errorf("mlbapi: invalid meta type %q; available types: %s", metaType, strings.Join(valid, ", "))
	}

	definition, ok := defaultEndpoints["meta"]
	if !ok {
		return nil, fmt.Errorf("mlbapi: meta endpoint definition not found")
	}

	requestURL, err := c.buildURL(definition, Params{"type": metaType}, false)
	if err != nil {
		return nil, err
	}

	var payload any
	if err := c.fetchJSON(ctx, "meta", requestURL, &payload); err != nil {
		return nil, err
	}

	switch typed := payload.(type) {
	case map[string]any:
		return typed, nil
	case []any:
		// The live meta endpoint returns a top-level array for many types.
		// Wrap it under the requested type so callers get a stable object shape.
		return JSON{metaType: typed}, nil
	default:
		return JSON{metaType: typed}, nil
	}
}

func playerLookupFromJSON(item JSON) PlayerLookup {
	currentTeam := teamFromJSON(asMap(item["currentTeam"]))
	return PlayerLookup{
		ID:            intValue(item["id"]),
		FullName:      stringValue(item["fullName"]),
		FirstName:     stringValue(item["firstName"]),
		LastName:      stringValue(item["lastName"]),
		PrimaryNumber: stringValue(item["primaryNumber"]),
		CurrentTeam: TeamReference{
			ID:           currentTeam.ID,
			Name:         currentTeam.Name,
			Abbreviation: currentTeam.Abbreviation,
		},
		PrimaryPosition: Position{
			Code:         stringValue(nestedValue(item, "primaryPosition", "code")),
			Name:         stringValue(nestedValue(item, "primaryPosition", "name")),
			Type:         stringValue(nestedValue(item, "primaryPosition", "type")),
			Abbreviation: stringValue(nestedValue(item, "primaryPosition", "abbreviation")),
		},
		Code:          stringValue(item["code"]),
		Abbreviation:  stringValue(item["abbreviation"]),
		UseName:       stringValue(item["useName"]),
		BoxscoreName:  stringValue(item["boxscoreName"]),
		Nickname:      stringValue(item["nickName"]),
		MLBDebutDate:  stringValue(item["mlbDebutDate"]),
		NameFirstLast: stringValue(item["nameFirstLast"]),
		FirstLastName: stringValue(item["firstLastName"]),
		LastFirstName: stringValue(item["lastFirstName"]),
		LastInitName:  stringValue(item["lastInitName"]),
		InitLastName:  stringValue(item["initLastName"]),
		FullFMLName:   stringValue(item["fullFMLName"]),
		FullLFMName:   stringValue(item["fullLFMName"]),
		NameSlug:      stringValue(item["nameSlug"]),
	}
}

func playerLookupMatches(player PlayerLookup, terms []string) bool {
	return containsAllTerms([]string{
		player.FullName,
		player.FirstName,
		player.LastName,
		player.UseName,
		player.BoxscoreName,
		player.Nickname,
		player.NameFirstLast,
		player.FirstLastName,
		player.LastFirstName,
		player.LastInitName,
		player.InitLastName,
		player.FullFMLName,
		player.FullLFMName,
		player.NameSlug,
		player.CurrentTeam.Name,
		player.CurrentTeam.Abbreviation,
		player.PrimaryPosition.Code,
		player.PrimaryPosition.Name,
		player.PrimaryPosition.Type,
		player.PrimaryPosition.Abbreviation,
		player.PrimaryNumber,
		player.Code,
		player.Abbreviation,
	}, terms)
}

func teamFromJSON(item JSON) Team {
	team := Team{
		ID:            TeamID(intValue(item["id"])),
		Name:          stringValue(item["name"]),
		TeamCode:      stringValue(item["teamCode"]),
		FileCode:      stringValue(item["fileCode"]),
		Abbreviation:  stringValue(item["abbreviation"]),
		TeamName:      stringValue(item["teamName"]),
		LocationName:  stringValue(item["locationName"]),
		ShortName:     stringValue(item["shortName"]),
		FranchiseName: stringValue(item["franchiseName"]),
		ClubName:      stringValue(item["clubName"]),
	}

	if known, ok := TeamByID(team.ID); ok {
		if team.Name == "" {
			team.Name = known.Name
		}
		if team.TeamCode == "" {
			team.TeamCode = known.TeamCode
		}
		if team.FileCode == "" {
			team.FileCode = known.FileCode
		}
		if team.Abbreviation == "" {
			team.Abbreviation = known.Abbreviation
		}
		if team.TeamName == "" {
			team.TeamName = known.TeamName
		}
		if team.LocationName == "" {
			team.LocationName = known.LocationName
		}
		if team.ShortName == "" {
			team.ShortName = known.ShortName
		}
		if team.FranchiseName == "" {
			team.FranchiseName = known.FranchiseName
		}
		if team.ClubName == "" {
			team.ClubName = known.ClubName
		}
	}

	return team
}

func teamMatches(team Team, target string) bool {
	return containsAllTerms([]string{
		team.Name,
		team.TeamCode,
		team.FileCode,
		team.Abbreviation,
		team.TeamName,
		team.LocationName,
		team.ShortName,
		team.FranchiseName,
		team.ClubName,
	}, []string{target})
}

func containsAllTerms(values, terms []string) bool {
	for _, term := range terms {
		found := false
		for _, value := range values {
			if strings.Contains(strings.ToLower(value), term) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func defaultString(item JSON, key, fallback string) string {
	if value := stringValue(item[key]); value != "" {
		return value
	}
	return fallback
}
