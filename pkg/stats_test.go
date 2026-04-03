package mlbapi

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestPlayerStatDataAndFormatting(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(r *http.Request) (int, string) {
		if r.URL.Path != "/api/v1/people/660271" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if !strings.Contains(r.URL.Query().Get("hydrate"), "season=2024") {
			t.Fatalf("hydrate missing requested season: %s", r.URL.Query().Get("hydrate"))
		}

		return http.StatusOK, `{
			"people":[
				{
					"id":660271,
					"useName":"Shohei",
					"lastName":"Ohtani",
					"active":true,
					"currentTeam":{"name":"Los Angeles Dodgers"},
					"primaryPosition":{"abbreviation":"DH"},
					"nickName":"Shotime",
					"lastPlayedDate":"2024-10-10",
					"mlbDebutDate":"2018-03-29",
					"batSide":{"description":"Left"},
					"pitchHand":{"description":"Right"},
					"stats":[
						{
							"type":{"displayName":"Season"},
							"group":{"displayName":"Hitting"},
							"splits":[
								{
									"season":"2024",
									"stat":{
										"avg":".310",
										"homeRuns":54,
										"position":{"abbreviation":"DH"}
									}
								}
							]
						}
					]
				}
			]
		}`
	})

	player, err := client.PlayerStatData(context.Background(), 660271, PlayerStatOptions{
		Type:    "season",
		Group:   "[hitting]",
		SportID: 1,
		Season:  2024,
	})
	if err != nil {
		t.Fatalf("PlayerStatData returned error: %v", err)
	}
	if player.FirstName != "Shohei" || len(player.Stats) != 1 {
		t.Fatalf("unexpected player summary: %+v", player)
	}

	rendered, err := client.PlayerStats(context.Background(), 660271, PlayerStatOptions{
		Type:    "season",
		Group:   "[hitting]",
		SportID: 1,
		Season:  2024,
	})
	if err != nil {
		t.Fatalf("PlayerStats returned error: %v", err)
	}
	if !strings.Contains(rendered, `Shohei "Shotime" Ohtani, DH (2018-)`) {
		t.Fatalf("formatted stats missing player header: %s", rendered)
	}
	if !strings.Contains(rendered, "homeRuns: 54") || !strings.Contains(rendered, "avg: .310") {
		t.Fatalf("formatted stats missing stat lines: %s", rendered)
	}
}

func TestLookupPlayerUsesLatestSeasonAndLookupTeamFilters(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(r *http.Request) (int, string) {
		switch r.URL.Path {
		case "/api/v1/seasons/all":
			return http.StatusOK, `{
				"seasons":[
					{"seasonId":"2024","seasonEndDate":"2024-11-01"},
					{"seasonId":"2999","seasonEndDate":"2999-11-01"}
				]
			}`
		case "/api/v1/sports/1/players":
			if got := r.URL.Query().Get("season"); got != "2999" {
				t.Fatalf("lookup player did not use latest season: %s", got)
			}
			return http.StatusOK, `{
				"people":[
					{"id":1,"fullName":"Mookie Betts","firstName":"Mookie","lastName":"Betts","useName":"Mookie","currentTeam":{"id":119},"primaryPosition":{"abbreviation":"RF"}},
					{"id":2,"fullName":"Freddie Freeman","firstName":"Freddie","lastName":"Freeman","useName":"Freddie","currentTeam":{"id":119},"primaryPosition":{"abbreviation":"1B"}}
				]
			}`
		case "/api/v1/teams":
			if got := r.URL.Query().Get("season"); got != "2024" {
				t.Fatalf("lookup team used unexpected season: %s", got)
			}
			return http.StatusOK, `{
				"teams":[
					{"id":119,"name":"Los Angeles Dodgers","teamCode":"lan","fileCode":"la","abbreviation":"LAD","teamName":"Dodgers","locationName":"Los Angeles","shortName":"LA Dodgers","franchiseName":"Los Angeles","clubName":"Dodgers"},
					{"id":147,"name":"New York Yankees","teamCode":"nya","fileCode":"nyy","abbreviation":"NYY","teamName":"Yankees","locationName":"Bronx","shortName":"NY Yankees","franchiseName":"New York","clubName":"Yankees"}
				]
			}`
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
			return 0, ""
		}
	})

	players, err := client.LookupPlayer(context.Background(), "mookie", LookupPlayerOptions{})
	if err != nil {
		t.Fatalf("LookupPlayer returned error: %v", err)
	}
	if len(players) != 1 || players[0].FullName != "Mookie Betts" {
		t.Fatalf("unexpected player results: %+v", players)
	}
	if players[0].CurrentTeam.ID != TeamLosAngelesDodgers {
		t.Fatalf("expected current team ID to be Dodgers, got %+v", players[0].CurrentTeam)
	}

	teams, err := client.LookupTeam(context.Background(), "dodg", LookupTeamOptions{
		Season:   2024,
		SportIDs: "1",
	})
	if err != nil {
		t.Fatalf("LookupTeam returned error: %v", err)
	}
	if len(teams) != 1 || teams[0].Name != "Los Angeles Dodgers" {
		t.Fatalf("unexpected team results: %+v", teams)
	}
	if teams[0].ID != TeamLosAngelesDodgers || teams[0].Abbreviation != "LAD" {
		t.Fatalf("unexpected typed team result: %+v", teams[0])
	}
}

func TestStandingsDataAndFormatting(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(r *http.Request) (int, string) {
		if r.URL.Path != "/api/v1/standings" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("season"); got != "2024" {
			t.Fatalf("unexpected season query: %s", got)
		}

		return http.StatusOK, `{
			"records":[
				{
					"teamRecords":[
						{
							"team":{"id":119,"name":"Los Angeles Dodgers","division":{"id":203,"name":"NL West","abbreviation":"W"}},
							"divisionRank":"1",
							"wins":100,
							"losses":62,
							"gamesBack":"-",
							"wildCardRank":"-",
							"wildCardGamesBack":"-",
							"wildCardEliminationNumber":"-",
							"eliminationNumber":"E",
							"leagueRank":"1",
							"sportRank":"1"
						},
						{
							"team":{"id":135,"name":"San Diego Padres","division":{"id":203,"name":"NL West","abbreviation":"W"}},
							"divisionRank":"2",
							"wins":90,
							"losses":72,
							"gamesBack":"10.0",
							"wildCardRank":"1",
							"wildCardGamesBack":"-",
							"wildCardEliminationNumber":"-",
							"eliminationNumber":"5",
							"leagueRank":"4",
							"sportRank":"4"
						},
						{
							"team":{"id":111,"name":"Boston Red Sox","division":{"id":201,"name":"AL East","abbreviation":"E"}},
							"divisionRank":"3",
							"wins":80,
							"losses":82,
							"gamesBack":"20.0",
							"wildCardRank":"6",
							"wildCardGamesBack":"8.0",
							"wildCardEliminationNumber":"3",
							"eliminationNumber":"7",
							"leagueRank":"8",
							"sportRank":"15"
						}
					]
				}
			]
		}`
	})

	divisions, err := client.StandingsData(context.Background(), StandingsOptions{
		Season:   2024,
		Division: "W",
	})
	if err != nil {
		t.Fatalf("StandingsData returned error: %v", err)
	}
	if len(divisions) != 1 || len(divisions[0].Teams) != 2 {
		t.Fatalf("unexpected standings divisions: %+v", divisions)
	}
	if divisions[0].DivisionName != "NL West" {
		t.Fatalf("unexpected division name: %+v", divisions[0])
	}
	if divisions[0].Teams[0].TeamID != TeamLosAngelesDodgers {
		t.Fatalf("unexpected typed team ID: %+v", divisions[0].Teams[0])
	}

	rendered, err := client.Standings(context.Background(), StandingsOptions{
		Season:   2024,
		Division: "W",
	})
	if err != nil {
		t.Fatalf("Standings returned error: %v", err)
	}
	if !strings.Contains(rendered, "NL West") || !strings.Contains(rendered, "Los Angeles Dodgers") {
		t.Fatalf("rendered standings missing expected content: %s", rendered)
	}
}

func TestMetaValidationAndNotes(t *testing.T) {
	t.Parallel()

	if _, err := DefaultClient.Meta(context.Background(), "not-a-real-meta-type"); err == nil {
		t.Fatal("expected invalid meta type error")
	}

	notes, err := Notes("schedule")
	if err != nil {
		t.Fatalf("Notes returned error: %v", err)
	}
	if !strings.Contains(notes, "All query parameters") {
		t.Fatalf("notes missing query parameter section: %s", notes)
	}
	if !strings.Contains(notes, "hydrate") {
		t.Fatalf("notes missing hydrate guidance: %s", notes)
	}
}

func TestMetaSuccessAndRosterFormatting(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(r *http.Request) (int, string) {
		switch r.URL.Path {
		case "/api/v1/positions":
			return http.StatusOK, `{
				"positions":[
					{"code":"1","name":"Pitcher","abbreviation":"P"}
				]
			}`
		case "/api/v1/teams/119/roster":
			if got := r.URL.Query().Get("rosterType"); got != "active" {
				t.Fatalf("unexpected rosterType query: %s", got)
			}
			return http.StatusOK, `{
				"roster":[
					{"jerseyNumber":"17","position":{"abbreviation":"DH"},"person":{"fullName":"Shohei Ohtani"}},
					{"jerseyNumber":"50","position":{"abbreviation":"SS"},"person":{"fullName":"Mookie Betts"}}
				]
			}`
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
			return 0, ""
		}
	})

	meta, err := client.Meta(context.Background(), "positions")
	if err != nil {
		t.Fatalf("Meta returned error: %v", err)
	}
	if len(nestedSlice(meta, "positions")) != 1 {
		t.Fatalf("unexpected meta payload: %+v", meta)
	}

	roster, err := client.Roster(context.Background(), TeamLosAngelesDodgers, RosterOptions{})
	if err != nil {
		t.Fatalf("Roster returned error: %v", err)
	}
	if !strings.Contains(roster, "#17  DH  Shohei Ohtani") {
		t.Fatalf("roster missing expected first row: %q", roster)
	}
	if !strings.Contains(roster, "#50  SS  Mookie Betts") {
		t.Fatalf("roster missing expected second row: %q", roster)
	}
}
