package mlbapi

import (
	"context"
	"fmt"
	"sort"
	"strings"
)

func Schedule(ctx context.Context, opts ScheduleOptions) ([]ScheduleGame, error) {
	return DefaultClient.Schedule(ctx, opts)
}

func (c *Client) Schedule(ctx context.Context, opts ScheduleOptions) ([]ScheduleGame, error) {
	if opts.EndDate != "" && opts.StartDate == "" {
		opts.Date = opts.EndDate
		opts.EndDate = ""
	}
	if opts.StartDate != "" && opts.EndDate == "" {
		opts.Date = opts.StartDate
		opts.StartDate = ""
	}

	params := Params{}
	switch {
	case opts.Date != "":
		params["date"] = opts.Date
	case opts.StartDate != "" && opts.EndDate != "":
		params["startDate"] = opts.StartDate
		params["endDate"] = opts.EndDate
	}

	if opts.TeamID != 0 {
		params["teamId"] = opts.TeamID
	}
	if opts.OpponentID != 0 {
		params["opponentId"] = opts.OpponentID
	}
	if opts.GameID != "" {
		params["gamePks"] = opts.GameID
	}
	if opts.LeagueID != "" {
		params["leagueId"] = opts.LeagueID
	}
	if opts.Season != 0 {
		params["season"] = opts.Season
	}

	sportID := opts.SportID
	if sportID == 0 {
		sportID = 1
	}

	hydrate := "decisions,probablePitcher(note),linescore,broadcasts,game(content(media(epg)))"
	if !opts.DisableSeriesStatus {
		switch {
		case opts.Date == "2014-03-11":
		case opts.StartDate != "" && opts.EndDate != "" && opts.StartDate <= "2014-03-11" && "2014-03-11" <= opts.EndDate:
		default:
			hydrate += ",seriesStatus"
		}
	}

	params["sportId"] = sportID
	params["hydrate"] = hydrate

	response, err := c.CallSchedule(ctx, params)
	if err != nil {
		return nil, err
	}

	if intValue(response["totalItems"]) == 0 {
		return []ScheduleGame{}, nil
	}

	games := make([]ScheduleGame, 0)
	for _, dateValue := range nestedSlice(response, "dates") {
		dateMap := asMap(dateValue)
		gameDate := stringValue(dateMap["date"])
		for _, gameValue := range asSlice(dateMap["games"]) {
			game := asMap(gameValue)
			awayTeam := nestedMap(game, "teams", "away", "team")
			homeTeam := nestedMap(game, "teams", "home", "team")
			away := nestedMap(game, "teams", "away")
			home := nestedMap(game, "teams", "home")
			linescore := asMap(game["linescore"])
			venue := asMap(game["venue"])
			content := asMap(game["content"])
			media := asMap(content["media"])
			seriesStatus := asMap(game["seriesStatus"])
			decisions := asMap(game["decisions"])

			nationalBroadcasts := make([]string, 0)
			for _, broadcastValue := range asSlice(game["broadcasts"]) {
				broadcast := asMap(broadcastValue)
				if boolValue(broadcast["isNational"]) {
					nationalBroadcasts = append(nationalBroadcasts, stringValue(broadcast["name"]))
				}
			}
			nationalBroadcasts = uniqueStrings(nationalBroadcasts)
			sort.Strings(nationalBroadcasts)
			if boolValue(media["freeGame"]) {
				nationalBroadcasts = append(nationalBroadcasts, "MLB.tv Free Game")
			}

			info := ScheduleGame{
				GameID:              intValue(game["gamePk"]),
				GameDatetime:        stringValue(game["gameDate"]),
				GameDate:            gameDate,
				GameType:            stringValue(game["gameType"]),
				Status:              nestedString(game, "status", "detailedState"),
				AwayName:            stringValue(awayTeam["name"]),
				HomeName:            stringValue(homeTeam["name"]),
				AwayID:              TeamID(intValue(awayTeam["id"])),
				HomeID:              TeamID(intValue(homeTeam["id"])),
				DoubleHeader:        stringValue(game["doubleHeader"]),
				GameNumber:          intValue(game["gameNumber"]),
				HomeProbablePitcher: nestedString(game, "teams", "home", "probablePitcher", "fullName"),
				AwayProbablePitcher: nestedString(game, "teams", "away", "probablePitcher", "fullName"),
				HomePitcherNote:     nestedString(game, "teams", "home", "probablePitcher", "note"),
				AwayPitcherNote:     nestedString(game, "teams", "away", "probablePitcher", "note"),
				AwayScore:           intValue(away["score"]),
				HomeScore:           intValue(home["score"]),
				CurrentInning:       intValue(linescore["currentInning"]),
				InningState:         stringValue(linescore["inningState"]),
				VenueID:             intValue(venue["id"]),
				VenueName:           stringValue(venue["name"]),
				NationalBroadcasts:  nationalBroadcasts,
				SeriesStatus:        stringValue(seriesStatus["result"]),
			}

			switch info.Status {
			case "Final", "Game Over":
				if boolValue(game["isTie"]) {
					info.WinningTeam = "Tie"
					info.LosingTeam = "Tie"
				} else if boolValue(away["isWinner"]) {
					info.WinningTeam = info.AwayName
					info.LosingTeam = info.HomeName
				} else {
					info.WinningTeam = info.HomeName
					info.LosingTeam = info.AwayName
				}
				info.WinningPitcher = nestedString(decisions, "winner", "fullName")
				info.LosingPitcher = nestedString(decisions, "loser", "fullName")
				info.SavePitcher = nestedString(decisions, "save", "fullName")
				info.Summary = fmt.Sprintf("%s - %s (%d) @ %s (%d) (%s)", gameDate, info.AwayName, info.AwayScore, info.HomeName, info.HomeScore, info.Status)
			case "In Progress":
				info.Summary = fmt.Sprintf("%s - %s (%d) @ %s (%d) (%s of the %s)", gameDate, info.AwayName, info.AwayScore, info.HomeName, info.HomeScore, info.InningState, nestedString(game, "linescore", "currentInningOrdinal"))
			default:
				info.Summary = fmt.Sprintf("%s - %s @ %s (%s)", gameDate, info.AwayName, info.HomeName, info.Status)
			}

			games = append(games, info)
		}
	}

	return games, nil
}

func LastGame(ctx context.Context, teamID TeamID) (int, bool, error) {
	return DefaultClient.LastGame(ctx, teamID)
}

func (c *Client) LastGame(ctx context.Context, teamID TeamID) (int, bool, error) {
	response, err := c.CallTeam(ctx, Params{
		"teamId":  teamID,
		"hydrate": "previousSchedule",
		"fields":  "teams,team,id,previousGameSchedule,dates,date,games,gamePk,gameDate,status,abstractGameCode",
	})
	if err != nil {
		return 0, false, err
	}

	teams := nestedSlice(response, "teams")
	if len(teams) == 0 {
		return 0, false, nil
	}

	dates := nestedSlice(asMap(teams[0]), "previousGameSchedule", "dates")
	lastCompleted := 0
	found := false
	for _, dateValue := range dates {
		for _, gameValue := range asSlice(asMap(dateValue)["games"]) {
			game := asMap(gameValue)
			if nestedString(game, "status", "abstractGameCode") == "F" {
				lastCompleted = intValue(game["gamePk"])
				found = true
			}
		}
	}

	return lastCompleted, found, nil
}

func NextGame(ctx context.Context, teamID TeamID) (int, bool, error) {
	return DefaultClient.NextGame(ctx, teamID)
}

func (c *Client) NextGame(ctx context.Context, teamID TeamID) (int, bool, error) {
	response, err := c.CallTeam(ctx, Params{
		"teamId":  teamID,
		"hydrate": "nextSchedule",
		"fields":  "teams,team,id,nextGameSchedule,dates,date,games,gamePk,gameDate,status,abstractGameCode",
	})
	if err != nil {
		return 0, false, err
	}

	teams := nestedSlice(response, "teams")
	if len(teams) == 0 {
		return 0, false, nil
	}

	for _, dateValue := range nestedSlice(asMap(teams[0]), "nextGameSchedule", "dates") {
		for _, gameValue := range asSlice(asMap(dateValue)["games"]) {
			game := asMap(gameValue)
			if nestedString(game, "status", "abstractGameCode") == "P" {
				return intValue(game["gamePk"]), true, nil
			}
		}
	}

	return 0, false, nil
}

func GameScoringPlays(ctx context.Context, gamePK int) (string, error) {
	return DefaultClient.GameScoringPlays(ctx, gamePK)
}

func (c *Client) GameScoringPlays(ctx context.Context, gamePK int) (string, error) {
	summary, err := c.GameScoringPlayData(ctx, gamePK)
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	for i, play := range summary.Plays {
		if i > 0 {
			builder.WriteString("\n\n")
		}
		builder.WriteString(nestedString(play, "result", "description"))
		builder.WriteByte('\n')
		builder.WriteString(fmt.Sprintf("%s %s - %s: %s, %s: %s",
			titleCaseFirst(nestedString(play, "about", "halfInning")),
			stringValue(nestedValue(play, "about", "inning")),
			stringValue(summary.Away["name"]),
			stringValue(nestedValue(play, "result", "awayScore")),
			stringValue(summary.Home["name"]),
			stringValue(nestedValue(play, "result", "homeScore")),
		))
	}

	return builder.String(), nil
}

func GameScoringPlayData(ctx context.Context, gamePK int) (ScoringPlayData, error) {
	return DefaultClient.GameScoringPlayData(ctx, gamePK)
}

func (c *Client) GameScoringPlayData(ctx context.Context, gamePK int) (ScoringPlayData, error) {
	response, err := c.CallGame(ctx, Params{
		"gamePk": gamePK,
		"fields": "gamePk,link,gameData,game,pk,teams,away,id,name,teamCode,fileCode,abbreviation,teamName,locationName,shortName,home,liveData,plays,allPlays,scoringPlays,scoringPlays,atBatIndex,result,description,awayScore,homeScore,about,halfInning,inning,endTime",
	})
	if err != nil {
		return ScoringPlayData{}, err
	}

	data := ScoringPlayData{
		Home: nestedMap(response, "gameData", "teams", "home"),
		Away: nestedMap(response, "gameData", "teams", "away"),
	}

	scoringPlays := nestedSlice(response, "liveData", "plays", "scoringPlays")
	if len(scoringPlays) == 0 {
		return data, nil
	}

	allPlays := nestedSlice(response, "liveData", "plays", "allPlays")
	ordered := make([]struct {
		EndTime string
		Play    JSON
	}, 0, len(scoringPlays))

	for _, scoringPlay := range scoringPlays {
		targetIndex := intValue(scoringPlay)
		for _, playValue := range allPlays {
			play := asMap(playValue)
			if intValue(play["atBatIndex"]) == targetIndex {
				ordered = append(ordered, struct {
					EndTime string
					Play    JSON
				}{
					EndTime: nestedString(play, "about", "endTime"),
					Play:    play,
				})
				break
			}
		}
	}

	sort.SliceStable(ordered, func(i, j int) bool {
		return ordered[i].EndTime < ordered[j].EndTime
	})

	data.Plays = make([]JSON, 0, len(ordered))
	for _, item := range ordered {
		data.Plays = append(data.Plays, item.Play)
	}

	return data, nil
}

func GameHighlights(ctx context.Context, gamePK int) (string, error) {
	return DefaultClient.GameHighlights(ctx, gamePK)
}

func (c *Client) GameHighlights(ctx context.Context, gamePK int) (string, error) {
	highlights, err := c.GameHighlightData(ctx, gamePK)
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	for i, item := range highlights {
		if i > 0 {
			builder.WriteString("\n\n")
		}
		builder.WriteString(fmt.Sprintf("%s (%s)\n%s\n%s",
			firstNonEmpty(item.Title, item.Headline),
			item.Duration,
			item.Description,
			playbackURL(item.Playbacks),
		))
	}

	return builder.String(), nil
}

func GameHighlightData(ctx context.Context, gamePK int) ([]HighlightItem, error) {
	return DefaultClient.GameHighlightData(ctx, gamePK)
}

func (c *Client) GameHighlightData(ctx context.Context, gamePK int) ([]HighlightItem, error) {
	response, err := c.CallSchedule(ctx, Params{
		"sportId": 1,
		"gamePk":  gamePK,
		"hydrate": "game(content(highlights(highlights)))",
		"fields":  "dates,date,games,gamePk,content,highlights,items,date,headline,type,value,title,description,duration,playbacks,name,url",
	})
	if err != nil {
		return nil, err
	}

	dates := nestedSlice(response, "dates")
	if len(dates) == 0 {
		return []HighlightItem{}, nil
	}

	games := asSlice(asMap(dates[0])["games"])
	if len(games) == 0 {
		return []HighlightItem{}, nil
	}

	items := nestedSlice(asMap(games[0]), "content", "highlights", "highlights", "items")
	if len(items) == 0 {
		return []HighlightItem{}, nil
	}

	highlights := make([]HighlightItem, 0, len(items))
	for _, itemValue := range items {
		item := asMap(itemValue)
		if stringValue(item["type"]) != "video" {
			continue
		}

		playbacks := make([]Playback, 0)
		for _, playbackValue := range asSlice(item["playbacks"]) {
			playback := asMap(playbackValue)
			playbacks = append(playbacks, Playback{
				Name: stringValue(playback["name"]),
				URL:  stringValue(playback["url"]),
			})
		}

		highlights = append(highlights, HighlightItem{
			Date:        stringValue(item["date"]),
			Title:       stringValue(item["title"]),
			Headline:    stringValue(item["headline"]),
			Description: stringValue(item["description"]),
			Duration:    stringValue(item["duration"]),
			Playbacks:   playbacks,
			Raw:         item,
		})
	}

	sort.SliceStable(highlights, func(i, j int) bool {
		return highlights[i].Date < highlights[j].Date
	})

	return highlights, nil
}

func nestedString(root JSON, path ...string) string {
	return stringValue(nestedValue(root, path...))
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func playbackURL(playbacks []Playback) string {
	for _, playback := range playbacks {
		if playback.Name == "mp4Avc" {
			return playback.URL
		}
	}
	for _, playback := range playbacks {
		if playback.Name == "FLASH_2500K_1280X720" {
			return playback.URL
		}
	}
	return "Link not found"
}

func titleCaseFirst(value string) string {
	if value == "" {
		return ""
	}
	return strings.ToUpper(value[:1]) + value[1:]
}
