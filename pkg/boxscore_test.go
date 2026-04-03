package mlbapi

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestLinescoreFormatsScoreboard(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(r *http.Request) (int, string) {
		if r.URL.Path != "/api/v1.1/game/123456/feed/live" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		return http.StatusOK, `{
			"gameData":{
				"status":{"abstractGameState":"Final"},
				"teams":{
					"away":{"teamName":"Away Club"},
					"home":{"teamName":"Home Club"}
				}
			},
			"liveData":{
				"linescore":{
					"innings":[
						{"num":1,"away":{"runs":2},"home":{"runs":0}},
						{"num":2,"away":{"runs":1},"home":{"runs":3}}
					],
					"teams":{
						"away":{"runs":3,"hits":8,"errors":0},
						"home":{"runs":3,"hits":7,"errors":1}
					}
				}
			}
		}`
	})

	linescore, err := client.Linescore(context.Background(), 123456, "")
	if err != nil {
		t.Fatalf("Linescore returned error: %v", err)
	}

	lines := strings.Split(linescore, "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d: %q", len(lines), linescore)
	}
	if !strings.Contains(linescore, "Final") || !strings.Contains(linescore, "Away Club") || !strings.Contains(linescore, "Home Club") {
		t.Fatalf("linescore missing expected labels: %s", linescore)
	}
	if !strings.Contains(linescore, "3   8   0") || !strings.Contains(linescore, "3   7   1") {
		t.Fatalf("linescore missing totals: %s", linescore)
	}
}

func TestBoxscoreDataParsesAndRenders(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(r *http.Request) (int, string) {
		if r.URL.Path != "/api/v1.1/game/123456/feed/live" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		return http.StatusOK, `{
			"gameData":{
				"game":{"id":"2024_04_03_awaymlb_homemlb_1"},
				"teams":{
					"away":{"teamName":"Away Club"},
					"home":{"teamName":"Home Club"}
				},
				"players":{
					"ID1":{"boxscoreName":"A. Batter"},
					"ID2":{"boxscoreName":"H. Batter"},
					"ID3":{"boxscoreName":"A. Pitcher"},
					"ID4":{"boxscoreName":"H. Pitcher"}
				}
			},
			"liveData":{
				"boxscore":{
					"info":[
						{"label":"Weather","value":"Sunny"}
					],
					"teams":{
						"away":{
							"batters":[1],
							"pitchers":[3],
							"players":{
								"ID1":{
									"battingOrder":"100",
									"position":{"abbreviation":"CF"},
									"stats":{"batting":{"note":"HR","atBats":4,"runs":2,"hits":2,"doubles":1,"triples":0,"homeRuns":1,"rbi":3,"stolenBases":0,"baseOnBalls":1,"strikeOuts":1,"leftOnBase":2}},
									"seasonStats":{"batting":{"avg":".300","ops":".900","obp":".360","slg":".540"}}
								},
								"ID3":{
									"stats":{"pitching":{"note":"W","inningsPitched":"6.0","hits":5,"runs":2,"earnedRuns":2,"baseOnBalls":1,"strikeOuts":7,"homeRuns":1,"pitchesThrown":92,"strikes":61}},
									"seasonStats":{"pitching":{"era":"2.50"}}
								}
							},
							"teamStats":{
								"batting":{"atBats":32,"runs":5,"hits":9,"homeRuns":1,"rbi":5,"baseOnBalls":3,"strikeOuts":8,"leftOnBase":7},
								"pitching":{"inningsPitched":"9.0","hits":7,"runs":3,"earnedRuns":3,"baseOnBalls":2,"strikeOuts":10,"homeRuns":1}
							},
							"note":[{"label":"2B","value":"A. Batter"}],
							"info":[{"title":"BATTING","fieldList":[{"label":"HR","value":"A. Batter (1)"}]}]
						},
						"home":{
							"batters":[2],
							"pitchers":[4],
							"players":{
								"ID2":{
									"battingOrder":"100",
									"position":{"abbreviation":"DH"},
									"stats":{"batting":{"note":"","atBats":4,"runs":1,"hits":1,"doubles":0,"triples":0,"homeRuns":0,"rbi":1,"stolenBases":0,"baseOnBalls":0,"strikeOuts":2,"leftOnBase":3}},
									"seasonStats":{"batting":{"avg":".250","ops":".710","obp":".320","slg":".390"}}
								},
								"ID4":{
									"stats":{"pitching":{"note":"L","inningsPitched":"5.0","hits":6,"runs":4,"earnedRuns":4,"baseOnBalls":2,"strikeOuts":5,"homeRuns":1,"numberOfPitches":88,"strikes":55}},
									"seasonStats":{"pitching":{"era":"3.80"}}
								}
							},
							"teamStats":{
								"batting":{"atBats":30,"runs":3,"hits":7,"homeRuns":1,"rbi":3,"baseOnBalls":2,"strikeOuts":9,"leftOnBase":6},
								"pitching":{"inningsPitched":"9.0","hits":9,"runs":5,"earnedRuns":5,"baseOnBalls":3,"strikeOuts":8,"homeRuns":1}
							},
							"note":[{"label":"HR","value":"H. Batter"}],
							"info":[{"title":"FIELDING","fieldList":[{"label":"E","value":"1"}]}]
						}
					}
				}
			}
		}`
	})

	data, err := client.BoxscoreData(context.Background(), 123456, "")
	if err != nil {
		t.Fatalf("BoxscoreData returned error: %v", err)
	}
	if len(data.AwayBatters) != 2 || len(data.HomeBatters) != 2 {
		t.Fatalf("unexpected batter lengths: away=%d home=%d", len(data.AwayBatters), len(data.HomeBatters))
	}
	if len(data.AwayPitchers) != 2 || len(data.HomePitchers) != 2 {
		t.Fatalf("unexpected pitcher lengths: away=%d home=%d", len(data.AwayPitchers), len(data.HomePitchers))
	}
	if data.AwayBattingNotes[0] != "2B-A. Batter" || data.HomeBattingNotes[0] != "HR-H. Batter" {
		t.Fatalf("unexpected batting notes: away=%v home=%v", data.AwayBattingNotes, data.HomeBattingNotes)
	}
	if data.AwayPitchers[1].P != "92" || data.HomePitchers[1].P != "88" {
		t.Fatalf("unexpected pitch counts: away=%+v home=%+v", data.AwayPitchers[1], data.HomePitchers[1])
	}

	boxscore, err := client.Boxscore(context.Background(), 123456, BoxscoreOptions{})
	if err != nil {
		t.Fatalf("Boxscore returned error: %v", err)
	}
	if !strings.Contains(boxscore, "Away Club Batters") || !strings.Contains(boxscore, "Home Club Pitchers") {
		t.Fatalf("boxscore missing team headers: %s", boxscore)
	}
	if !strings.Contains(boxscore, "Weather: Sunny") {
		t.Fatalf("boxscore missing game info: %s", boxscore)
	}
	if !strings.Contains(boxscore, "Totals") {
		t.Fatalf("boxscore missing totals row: %s", boxscore)
	}
}
