package mlbapi

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func TestBuildURLUsesDefaultsAndBoolPathParams(t *testing.T) {
	t.Parallel()

	client := NewClient(WithBaseURL("https://example.com/api/"))
	url, err := client.buildURL(defaultEndpoints["awards"], Params{
		"awardId":    "mvp",
		"recipients": true,
		"season":     2024,
	}, false)
	if err != nil {
		t.Fatalf("buildURL returned error: %v", err)
	}

	want := "https://example.com/api/v1/awards/mvp/recipients?season=2024"
	if url != want {
		t.Fatalf("unexpected url:\nwant: %s\n got: %s", want, url)
	}
}

func TestGetRequiresQueryParamSet(t *testing.T) {
	t.Parallel()

	client := NewClient()
	_, err := client.Get(context.Background(), "game_diff", Params{
		"gamePk": 1234,
	})
	if err == nil {
		t.Fatal("expected missing parameter error")
	}
	if !strings.Contains(err.Error(), "missing required query parameters") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestScheduleTransformsResponse(t *testing.T) {
	t.Parallel()

	client := NewClient(
		WithBaseURL("https://example.test/api"),
		WithHTTPClient(&http.Client{
			Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				if r.URL.Path != "/api/v1/schedule" {
					t.Fatalf("unexpected path: %s", r.URL.Path)
				}
				if got := r.URL.Query().Get("date"); got != "2024-04-03" {
					t.Fatalf("unexpected date query: %s", got)
				}
				if got := r.URL.Query().Get("sportId"); got != "1" {
					t.Fatalf("unexpected sportId query: %s", got)
				}

				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{"Content-Type": []string{"application/json"}},
					Body: io.NopCloser(strings.NewReader(`{
			"totalItems": 1,
			"dates": [
				{
					"date": "2024-04-03",
					"games": [
						{
							"gamePk": 123456,
							"gameDate": "2024-04-03T17:10:00Z",
							"gameType": "R",
							"doubleHeader": "N",
							"gameNumber": 1,
							"status": {"detailedState": "Final"},
							"teams": {
								"away": {
									"score": 5,
									"isWinner": true,
									"team": {"id": 111, "name": "Away Club"},
									"probablePitcher": {"fullName": "Away Starter", "note": "RHP"}
								},
								"home": {
									"score": 3,
									"team": {"id": 222, "name": "Home Club"},
									"probablePitcher": {"fullName": "Home Starter", "note": "LHP"}
								}
							},
							"venue": {"id": 10, "name": "Test Park"},
							"linescore": {"currentInning": 9, "inningState": "End"},
							"broadcasts": [
								{"name": "ESPN", "isNational": true},
								{"name": "ESPN", "isNational": true},
								{"name": "Local", "isNational": false}
							],
							"content": {"media": {"freeGame": true}},
							"seriesStatus": {"result": "Away Club leads 1-0"},
							"decisions": {
								"winner": {"fullName": "Winning Pitcher"},
								"loser": {"fullName": "Losing Pitcher"},
								"save": {"fullName": "Save Pitcher"}
							}
						}
					]
					}
				]
		}`)),
				}, nil
			}),
		}),
	)
	games, err := client.Schedule(context.Background(), ScheduleOptions{Date: "2024-04-03"})
	if err != nil {
		t.Fatalf("Schedule returned error: %v", err)
	}
	if len(games) != 1 {
		t.Fatalf("expected 1 game, got %d", len(games))
	}

	game := games[0]
	if game.GameID != 123456 {
		t.Fatalf("unexpected game id: %d", game.GameID)
	}
	if game.WinningTeam != "Away Club" || game.LosingTeam != "Home Club" {
		t.Fatalf("unexpected winner/loser: %+v", game)
	}
	if len(game.NationalBroadcasts) != 2 {
		t.Fatalf("expected 2 national broadcasts, got %v", game.NationalBroadcasts)
	}
	if game.NationalBroadcasts[0] != "ESPN" || game.NationalBroadcasts[1] != "MLB.tv Free Game" {
		t.Fatalf("unexpected broadcasts: %v", game.NationalBroadcasts)
	}
	if game.Summary != "2024-04-03 - Away Club (5) @ Home Club (3) (Final)" {
		t.Fatalf("unexpected summary: %s", game.Summary)
	}
}

func TestGameScoringPlayDataSortsByEndTime(t *testing.T) {
	t.Parallel()

	client := NewClient(
		WithBaseURL("https://example.test/api"),
		WithHTTPClient(&http.Client{
			Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{"Content-Type": []string{"application/json"}},
					Body: io.NopCloser(strings.NewReader(`{
			"gameData": {
				"teams": {
					"home": {"name": "Home Club"},
					"away": {"name": "Away Club"}
				}
			},
			"liveData": {
				"plays": {
					"scoringPlays": [2, 1],
					"allPlays": [
						{
							"atBatIndex": 1,
							"result": {"description": "First run", "awayScore": 1, "homeScore": 0},
							"about": {"halfInning": "top", "inning": 1, "endTime": "2024-04-03T17:20:00Z"}
						},
						{
							"atBatIndex": 2,
							"result": {"description": "Second run", "awayScore": 1, "homeScore": 1},
							"about": {"halfInning": "bottom", "inning": 1, "endTime": "2024-04-03T17:30:00Z"}
						}
					]
					}
				}
		}`)),
				}, nil
			}),
		}),
	)
	data, err := client.GameScoringPlayData(context.Background(), 123456)
	if err != nil {
		t.Fatalf("GameScoringPlayData returned error: %v", err)
	}
	if len(data.Plays) != 2 {
		t.Fatalf("expected 2 plays, got %d", len(data.Plays))
	}
	if nestedString(data.Plays[0], "result", "description") != "First run" {
		t.Fatalf("plays not sorted by endTime: %+v", data.Plays)
	}
}
