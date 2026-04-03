package mlbapi

import (
	"context"
	"fmt"
	"strings"
)

// Boxscore renders a text boxscore for a game.
func Boxscore(ctx context.Context, gamePK int, opts BoxscoreOptions) (string, error) {
	return DefaultClient.Boxscore(ctx, gamePK, opts)
}

// Boxscore renders a text boxscore for a game.
func (c *Client) Boxscore(ctx context.Context, gamePK int, opts BoxscoreOptions) (string, error) {
	data, err := c.BoxscoreData(ctx, gamePK, opts.Timecode)
	if err != nil {
		return "", err
	}

	const rowLen = 79
	const fullRowLen = rowLen*2 + 3

	var builder strings.Builder

	if !opts.SkipBattingBox {
		awayBatters := append([]BatterLine(nil), data.AwayBatters...)
		homeBatters := append([]BatterLine(nil), data.HomeBatters...)
		blank := BatterLine{}

		for len(awayBatters) > len(homeBatters) {
			homeBatters = append(homeBatters, blank)
		}
		for len(homeBatters) > len(awayBatters) {
			awayBatters = append(awayBatters, blank)
		}

		awayBatters = append(awayBatters, data.AwayBattingTotals)
		homeBatters = append(homeBatters, data.HomeBattingTotals)

		for i := range awayBatters {
			if i == 0 || i == len(awayBatters)-1 {
				builder.WriteString(strings.Repeat("-", rowLen))
				builder.WriteString(" | ")
				builder.WriteString(strings.Repeat("-", rowLen))
				builder.WriteByte('\n')
			}

			builder.WriteString(formatBatterDisplay(awayBatters[i]))
			builder.WriteString(" | ")
			builder.WriteString(formatBatterDisplay(homeBatters[i]))
			builder.WriteByte('\n')

			if i == 0 || i == len(awayBatters)-1 {
				builder.WriteString(strings.Repeat("-", rowLen))
				builder.WriteString(" | ")
				builder.WriteString(strings.Repeat("-", rowLen))
				builder.WriteByte('\n')
			}
		}

		awayNotes := append([]string(nil), data.AwayBattingNotes...)
		homeNotes := append([]string(nil), data.HomeBattingNotes...)
		for len(awayNotes) > len(homeNotes) {
			homeNotes = append(homeNotes, "")
		}
		for len(homeNotes) > len(awayNotes) {
			awayNotes = append(awayNotes, "")
		}

		for i := range awayNotes {
			builder.WriteString(left(awayNotes[i], rowLen))
			builder.WriteString(" | ")
			builder.WriteString(left(homeNotes[i], rowLen))
			builder.WriteByte('\n')
		}

		builder.WriteString(strings.Repeat(" ", rowLen))
		builder.WriteString(" | ")
		builder.WriteString(strings.Repeat(" ", rowLen))
		builder.WriteByte('\n')
	}

	awayInfo := buildInfoBoxLines(data.Away, !opts.SkipBattingInfo, !opts.SkipFieldingInfo, rowLen)
	homeInfo := buildInfoBoxLines(data.Home, !opts.SkipBattingInfo, !opts.SkipFieldingInfo, rowLen)
	if len(awayInfo) > 0 || len(homeInfo) > 0 {
		for len(awayInfo) > len(homeInfo) {
			homeInfo = append(homeInfo, "")
		}
		for len(homeInfo) > len(awayInfo) {
			awayInfo = append(awayInfo, "")
		}

		for i := range awayInfo {
			builder.WriteString(left(awayInfo[i], rowLen))
			builder.WriteString(" | ")
			builder.WriteString(left(homeInfo[i], rowLen))
			builder.WriteByte('\n')
			if i == len(awayInfo)-1 {
				builder.WriteString(strings.Repeat("-", rowLen))
				builder.WriteString(" | ")
				builder.WriteString(strings.Repeat("-", rowLen))
				builder.WriteByte('\n')
			}
		}
	}

	if !opts.SkipPitchingBox {
		awayPitchers := append([]PitcherLine(nil), data.AwayPitchers...)
		homePitchers := append([]PitcherLine(nil), data.HomePitchers...)
		blank := PitcherLine{}

		for len(awayPitchers) > len(homePitchers) {
			homePitchers = append(homePitchers, blank)
		}
		for len(homePitchers) > len(awayPitchers) {
			awayPitchers = append(awayPitchers, blank)
		}

		awayPitchers = append(awayPitchers, data.AwayPitchingTotals)
		homePitchers = append(homePitchers, data.HomePitchingTotals)

		for i := range awayPitchers {
			if i == 0 || i == len(awayPitchers)-1 {
				builder.WriteString(strings.Repeat("-", rowLen))
				builder.WriteString(" | ")
				builder.WriteString(strings.Repeat("-", rowLen))
				builder.WriteByte('\n')
			}

			builder.WriteString(formatPitcherDisplay(awayPitchers[i]))
			builder.WriteString(" | ")
			builder.WriteString(formatPitcherDisplay(homePitchers[i]))
			builder.WriteByte('\n')

			if i == 0 || i == len(awayPitchers)-1 {
				builder.WriteString(strings.Repeat("-", rowLen))
				builder.WriteString(" | ")
				builder.WriteString(strings.Repeat("-", rowLen))
				builder.WriteByte('\n')
			}
		}
	}

	if !opts.SkipGameInfo {
		gameInfoLines := formatGameInfoLines(data.GameBoxInfo, fullRowLen)
		for i, line := range gameInfoLines {
			builder.WriteString(left(line, fullRowLen))
			builder.WriteByte('\n')
			if i == len(gameInfoLines)-1 {
				builder.WriteString(strings.Repeat("-", fullRowLen))
				builder.WriteByte('\n')
			}
		}
	}

	return builder.String(), nil
}

// BoxscoreDataForGame returns parsed boxscore data for a game.
func BoxscoreDataForGame(ctx context.Context, gamePK int, timecode string) (BoxscoreData, error) {
	return DefaultClient.BoxscoreData(ctx, gamePK, timecode)
}

// BoxscoreData returns parsed boxscore data for a game.
func (c *Client) BoxscoreData(ctx context.Context, gamePK int, timecode string) (BoxscoreData, error) {
	params := Params{
		"gamePk": gamePK,
		"fields": "gameData,game,teams,teamName,shortName,teamStats,batting,atBats,runs,hits,doubles,triples,homeRuns,rbi,stolenBases,strikeOuts,baseOnBalls,leftOnBase,pitching,inningsPitched,earnedRuns,homeRuns,players,boxscoreName,liveData,boxscore,teams,players,id,fullName,allPositions,abbreviation,seasonStats,batting,avg,ops,obp,slg,era,pitchesThrown,numberOfPitches,strikes,battingOrder,info,title,fieldList,note,label,value,wins,losses,holds,blownSaves",
	}
	if timecode != "" {
		params["timecode"] = timecode
	}

	response, err := c.CallGame(ctx, params)
	if err != nil {
		return BoxscoreData{}, err
	}

	teamInfo := nestedMap(response, "gameData", "teams")
	playerInfo := nestedMap(response, "gameData", "players")
	away := nestedMap(response, "liveData", "boxscore", "teams", "away")
	home := nestedMap(response, "liveData", "boxscore", "teams", "home")

	data := BoxscoreData{
		GameID:     stringValue(nestedValue(response, "gameData", "game", "id")),
		TeamInfo:   teamInfo,
		PlayerInfo: playerInfo,
		Away:       away,
		Home:       home,
		AwayBatters: []BatterLine{{
			NameField: teamName(teamInfo, "away") + " Batters",
			AB:        "AB",
			R:         "R",
			H:         "H",
			Doubles:   "2B",
			Triples:   "3B",
			HR:        "HR",
			RBI:       "RBI",
			SB:        "SB",
			BB:        "BB",
			K:         "K",
			LOB:       "LOB",
			AVG:       "AVG",
			OPS:       "OPS",
			Name:      teamName(teamInfo, "away") + " Batters",
			OBP:       "OBP",
			SLG:       "SLG",
		}},
		HomeBatters: []BatterLine{{
			NameField: teamName(teamInfo, "home") + " Batters",
			AB:        "AB",
			R:         "R",
			H:         "H",
			Doubles:   "2B",
			Triples:   "3B",
			HR:        "HR",
			RBI:       "RBI",
			SB:        "SB",
			BB:        "BB",
			K:         "K",
			LOB:       "LOB",
			AVG:       "AVG",
			OPS:       "OPS",
			Name:      teamName(teamInfo, "home") + " Batters",
			OBP:       "OBP",
			SLG:       "SLG",
		}},
		AwayPitchers: []PitcherLine{{
			NameField: teamName(teamInfo, "away") + " Pitchers",
			IP:        "IP",
			H:         "H",
			R:         "R",
			ER:        "ER",
			BB:        "BB",
			K:         "K",
			HR:        "HR",
			ERA:       "ERA",
			P:         "P",
			S:         "S",
			Name:      teamName(teamInfo, "away") + " Pitchers",
		}},
		HomePitchers: []PitcherLine{{
			NameField: teamName(teamInfo, "home") + " Pitchers",
			IP:        "IP",
			H:         "H",
			R:         "R",
			ER:        "ER",
			BB:        "BB",
			K:         "K",
			HR:        "HR",
			ERA:       "ERA",
			P:         "P",
			S:         "S",
			Name:      teamName(teamInfo, "home") + " Pitchers",
		}},
	}

	for _, sideName := range []string{"away", "home"} {
		side := away
		if sideName == "home" {
			side = home
		}

		for _, batterIDValue := range asSlice(side["batters"]) {
			batterID := intValue(batterIDValue)
			playerKey := fmt.Sprintf("ID%d", batterID)
			sidePlayer := nestedMap(side, "players", playerKey)
			if stringValue(sidePlayer["battingOrder"]) == "" {
				continue
			}

			battingStats := nestedMap(sidePlayer, "stats", "batting")
			if len(battingStats) == 0 {
				continue
			}

			order := stringValue(sidePlayer["battingOrder"])
			nameField := "   "
			if strings.HasSuffix(order, "0") && len(order) > 0 {
				nameField = order[:1]
			}
			nameField += " " + stringValue(battingStats["note"])
			nameField += nestedString(playerInfo, playerKey, "boxscoreName") + "  " + nestedString(sidePlayer, "position", "abbreviation")

			line := BatterLine{
				NameField:    nameField,
				AB:           stringValue(battingStats["atBats"]),
				R:            stringValue(battingStats["runs"]),
				H:            stringValue(battingStats["hits"]),
				Doubles:      stringValue(battingStats["doubles"]),
				Triples:      stringValue(battingStats["triples"]),
				HR:           stringValue(battingStats["homeRuns"]),
				RBI:          stringValue(battingStats["rbi"]),
				SB:           stringValue(battingStats["stolenBases"]),
				BB:           stringValue(battingStats["baseOnBalls"]),
				K:            stringValue(battingStats["strikeOuts"]),
				LOB:          stringValue(battingStats["leftOnBase"]),
				AVG:          nestedString(sidePlayer, "seasonStats", "batting", "avg"),
				OPS:          nestedString(sidePlayer, "seasonStats", "batting", "ops"),
				PersonID:     batterID,
				BattingOrder: order,
				Substitution: order != "" && !strings.HasSuffix(order, "0"),
				Note:         stringValue(battingStats["note"]),
				Name:         nestedString(playerInfo, playerKey, "boxscoreName"),
				Position:     nestedString(sidePlayer, "position", "abbreviation"),
				OBP:          nestedString(sidePlayer, "seasonStats", "batting", "obp"),
				SLG:          nestedString(sidePlayer, "seasonStats", "batting", "slg"),
			}

			if sideName == "away" {
				data.AwayBatters = append(data.AwayBatters, line)
			} else {
				data.HomeBatters = append(data.HomeBatters, line)
			}
		}

		totals := BatterLine{
			NameField: "Totals",
			AB:        nestedString(side, "teamStats", "batting", "atBats"),
			R:         nestedString(side, "teamStats", "batting", "runs"),
			H:         nestedString(side, "teamStats", "batting", "hits"),
			HR:        nestedString(side, "teamStats", "batting", "homeRuns"),
			RBI:       nestedString(side, "teamStats", "batting", "rbi"),
			BB:        nestedString(side, "teamStats", "batting", "baseOnBalls"),
			K:         nestedString(side, "teamStats", "batting", "strikeOuts"),
			LOB:       nestedString(side, "teamStats", "batting", "leftOnBase"),
			Name:      "Totals",
		}

		notes := make([]string, 0)
		for _, noteValue := range asSlice(side["note"]) {
			note := asMap(noteValue)
			notes = append(notes, stringValue(note["label"])+"-"+stringValue(note["value"]))
		}

		for _, pitcherIDValue := range asSlice(side["pitchers"]) {
			pitcherID := intValue(pitcherIDValue)
			playerKey := fmt.Sprintf("ID%d", pitcherID)
			sidePlayer := nestedMap(side, "players", playerKey)
			if len(sidePlayer) == 0 {
				continue
			}

			pitchingStats := nestedMap(sidePlayer, "stats", "pitching")
			if len(pitchingStats) == 0 {
				continue
			}

			nameField := nestedString(playerInfo, playerKey, "boxscoreName")
			if note := stringValue(pitchingStats["note"]); note != "" {
				nameField += "  " + note
			}

			pitchesThrown := stringValue(pitchingStats["pitchesThrown"])
			if pitchesThrown == "" {
				pitchesThrown = stringValue(pitchingStats["numberOfPitches"])
			}

			line := PitcherLine{
				NameField: nameField,
				IP:        stringValue(pitchingStats["inningsPitched"]),
				H:         stringValue(pitchingStats["hits"]),
				R:         stringValue(pitchingStats["runs"]),
				ER:        stringValue(pitchingStats["earnedRuns"]),
				BB:        stringValue(pitchingStats["baseOnBalls"]),
				K:         stringValue(pitchingStats["strikeOuts"]),
				HR:        stringValue(pitchingStats["homeRuns"]),
				P:         pitchesThrown,
				S:         stringValue(pitchingStats["strikes"]),
				ERA:       nestedString(sidePlayer, "seasonStats", "pitching", "era"),
				Name:      nestedString(playerInfo, playerKey, "boxscoreName"),
				PersonID:  pitcherID,
				Note:      stringValue(pitchingStats["note"]),
			}

			if sideName == "away" {
				data.AwayPitchers = append(data.AwayPitchers, line)
			} else {
				data.HomePitchers = append(data.HomePitchers, line)
			}
		}

		pitchingTotals := PitcherLine{
			NameField: "Totals",
			IP:        nestedString(side, "teamStats", "pitching", "inningsPitched"),
			H:         nestedString(side, "teamStats", "pitching", "hits"),
			R:         nestedString(side, "teamStats", "pitching", "runs"),
			ER:        nestedString(side, "teamStats", "pitching", "earnedRuns"),
			BB:        nestedString(side, "teamStats", "pitching", "baseOnBalls"),
			K:         nestedString(side, "teamStats", "pitching", "strikeOuts"),
			HR:        nestedString(side, "teamStats", "pitching", "homeRuns"),
			Name:      "Totals",
		}

		if sideName == "away" {
			data.AwayBattingTotals = totals
			data.AwayBattingNotes = notes
			data.AwayPitchingTotals = pitchingTotals
		} else {
			data.HomeBattingTotals = totals
			data.HomeBattingNotes = notes
			data.HomePitchingTotals = pitchingTotals
		}
	}

	for _, infoValue := range nestedSlice(response, "liveData", "boxscore", "info") {
		info := asMap(infoValue)
		data.GameBoxInfo = append(data.GameBoxInfo, InfoField{
			Label: stringValue(info["label"]),
			Value: stringValue(info["value"]),
		})
	}

	return data, nil
}

// Linescore renders an inning-by-inning linescore for a game.
func Linescore(ctx context.Context, gamePK int, timecode string) (string, error) {
	return DefaultClient.Linescore(ctx, gamePK, timecode)
}

// Linescore renders an inning-by-inning linescore for a game.
func (c *Client) Linescore(ctx context.Context, gamePK int, timecode string) (string, error) {
	params := Params{
		"gamePk": gamePK,
		"fields": "gameData,teams,teamName,shortName,status,abstractGameState,liveData,linescore,innings,num,home,away,runs,hits,errors",
	}
	if timecode != "" {
		params["timecode"] = timecode
	}

	response, err := c.CallGame(ctx, params)
	if err != nil {
		return "", err
	}

	headerName := nestedString(response, "gameData", "status", "abstractGameState")
	awayName := nestedString(response, "gameData", "teams", "away", "teamName")
	homeName := nestedString(response, "gameData", "teams", "home", "teamName")

	headerRow := make([]string, 0)
	away := make([]string, 0)
	home := make([]string, 0)

	for _, inningValue := range nestedSlice(response, "liveData", "linescore", "innings") {
		inning := asMap(inningValue)
		headerRow = append(headerRow, stringValue(inning["num"]))
		away = append(away, stringValue(nestedValue(inning, "away", "runs")))
		home = append(home, stringValue(nestedValue(inning, "home", "runs")))
	}

	for len(headerRow) < 9 {
		headerRow = append(headerRow, fmt.Sprintf("%d", len(headerRow)+1))
		away = append(away, " ")
		home = append(home, " ")
	}

	headerRow = append(headerRow, "R", "H", "E")
	away = append(away,
		nestedString(response, "liveData", "linescore", "teams", "away", "runs"),
		nestedString(response, "liveData", "linescore", "teams", "away", "hits"),
		nestedString(response, "liveData", "linescore", "teams", "away", "errors"),
	)
	home = append(home,
		nestedString(response, "liveData", "linescore", "teams", "home", "runs"),
		nestedString(response, "liveData", "linescore", "teams", "home", "hits"),
		nestedString(response, "liveData", "linescore", "teams", "home", "errors"),
	)

	width := len(headerName)
	if len(awayName) > width {
		width = len(awayName)
	}
	if len(homeName) > width {
		width = len(homeName)
	}
	width++

	rows := []struct {
		Name   string
		Values []string
	}{
		{Name: headerName, Values: headerRow},
		{Name: awayName, Values: away},
		{Name: homeName, Values: home},
	}

	var builder strings.Builder
	for idx, row := range rows {
		builder.WriteString(left(row.Name, width))
		for _, value := range row.Values[:len(row.Values)-3] {
			builder.WriteString(center(value, 2))
		}
		for _, value := range row.Values[len(row.Values)-3:] {
			builder.WriteString(center(value, 4))
		}
		if idx < len(rows)-1 {
			builder.WriteByte('\n')
		}
	}

	return builder.String(), nil
}

func teamName(teamInfo JSON, side string) string {
	return nestedString(teamInfo, side, "teamName")
}

func formatBatterDisplay(line BatterLine) string {
	return left(line.NameField, 40) +
		" " + center(line.AB, 3) +
		" " + center(line.R, 3) +
		" " + center(line.H, 3) +
		" " + center(line.RBI, 3) +
		" " + center(line.BB, 3) +
		" " + center(line.K, 3) +
		" " + center(line.LOB, 3) +
		" " + center(line.AVG, 4) +
		" " + center(line.OPS, 5)
}

func formatPitcherDisplay(line PitcherLine) string {
	return left(line.NameField, 43) +
		" " + center(line.IP, 4) +
		" " + center(line.H, 3) +
		" " + center(line.R, 3) +
		" " + center(line.ER, 3) +
		" " + center(line.BB, 3) +
		" " + center(line.K, 3) +
		" " + center(line.HR, 3) +
		" " + center(line.ERA, 6)
}

func buildInfoBoxLines(side JSON, includeBatting, includeFielding bool, rowLen int) []string {
	lines := make([]string, 0)
	for _, infoType := range []string{"BATTING", "FIELDING"} {
		if infoType == "BATTING" && !includeBatting {
			continue
		}
		if infoType == "FIELDING" && !includeFielding {
			continue
		}

		startLen := len(lines)
		for _, infoValue := range asSlice(side["info"]) {
			info := asMap(infoValue)
			if stringValue(info["title"]) != infoType {
				continue
			}

			lines = append(lines, infoType)
			for _, fieldValue := range asSlice(info["fieldList"]) {
				field := asMap(fieldValue)
				text := stringValue(field["label"]) + ": " + stringValue(field["value"])
				if len(text) > rowLen {
					lines = append(lines, wrapText(text, rowLen)...)
				} else {
					lines = append(lines, text)
				}
			}
		}

		if infoType == "BATTING" && len(lines) > startLen {
			lines = append(lines, " ")
		}
	}

	return lines
}

func formatGameInfoLines(fields []InfoField, width int) []string {
	lines := make([]string, 0)
	for _, field := range fields {
		text := field.Label
		if field.Value != "" {
			text += ": " + field.Value
		}
		if len(text) > width {
			lines = append(lines, wrapText(text, width)...)
			continue
		}
		lines = append(lines, text)
	}
	return lines
}
