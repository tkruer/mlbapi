package mlbapi

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestScheduleOmitsSeriesStatusHydrateOnKnownBadDate(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(r *http.Request) (int, string) {
		if r.URL.Path != "/api/v1/schedule" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if strings.Contains(r.URL.Query().Get("hydrate"), "seriesStatus") {
			t.Fatalf("hydrate unexpectedly included seriesStatus: %s", r.URL.Query().Get("hydrate"))
		}
		return http.StatusOK, `{"totalItems":0,"dates":[]}`
	})

	games, err := client.Schedule(context.Background(), ScheduleOptions{Date: "2014-03-11"})
	if err != nil {
		t.Fatalf("Schedule returned error: %v", err)
	}
	if len(games) != 0 {
		t.Fatalf("expected no games, got %d", len(games))
	}
}

func TestLastGameAndNextGame(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(r *http.Request) (int, string) {
		if r.URL.Path != "/api/v1/teams/109" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		switch r.URL.Query().Get("hydrate") {
		case "previousSchedule":
			return http.StatusOK, `{
				"teams":[
					{
						"previousGameSchedule":{
							"dates":[
								{"games":[
									{"gamePk":1001,"status":{"abstractGameCode":"P"}},
									{"gamePk":1002,"status":{"abstractGameCode":"F"}}
								]},
								{"games":[
									{"gamePk":1003,"status":{"abstractGameCode":"F"}}
								]}
							]
						}
					}
				]
			}`
		case "nextSchedule":
			return http.StatusOK, `{
				"teams":[
					{
						"nextGameSchedule":{
							"dates":[
								{"games":[
									{"gamePk":2001,"status":{"abstractGameCode":"S"}},
									{"gamePk":2002,"status":{"abstractGameCode":"P"}}
								]},
								{"games":[
									{"gamePk":2003,"status":{"abstractGameCode":"P"}}
								]}
							]
						}
					}
				]
			}`
		default:
			t.Fatalf("unexpected hydrate query: %s", r.URL.RawQuery)
			return 0, ""
		}
	})

	lastGame, found, err := client.LastGame(context.Background(), 109)
	if err != nil {
		t.Fatalf("LastGame returned error: %v", err)
	}
	if !found || lastGame != 1003 {
		t.Fatalf("unexpected last game result: found=%v game=%d", found, lastGame)
	}

	nextGame, found, err := client.NextGame(context.Background(), 109)
	if err != nil {
		t.Fatalf("NextGame returned error: %v", err)
	}
	if !found || nextGame != 2002 {
		t.Fatalf("unexpected next game result: found=%v game=%d", found, nextGame)
	}
}

func TestGameHighlightsSortAndFormat(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(r *http.Request) (int, string) {
		if r.URL.Path != "/api/v1/schedule" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("gamePk"); got != "123456" {
			t.Fatalf("unexpected gamePk query: %s", got)
		}

		return http.StatusOK, `{
			"dates":[
				{
					"games":[
						{
							"content":{
								"highlights":{
									"highlights":{
										"items":[
											{
												"date":"2024-04-03T17:30:00Z",
												"type":"video",
												"headline":"Later highlight",
												"description":"Later play",
												"duration":"01:00",
												"playbacks":[{"name":"FLASH_2500K_1280X720","url":"https://video.example/later"}]
											},
											{
												"date":"2024-04-03T17:20:00Z",
												"type":"video",
												"title":"Earlier highlight",
												"description":"Earlier play",
												"duration":"00:30",
												"playbacks":[{"name":"mp4Avc","url":"https://video.example/earlier"}]
											},
											{
												"date":"2024-04-03T17:10:00Z",
												"type":"article",
												"title":"Ignore me"
											}
										]
									}
								}
							}
						}
					]
				}
			]
		}`
	})

	items, err := client.GameHighlightData(context.Background(), 123456)
	if err != nil {
		t.Fatalf("GameHighlightData returned error: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 highlight items, got %d", len(items))
	}
	if items[0].Title != "Earlier highlight" {
		t.Fatalf("expected earliest highlight first, got %+v", items)
	}

	rendered, err := client.GameHighlights(context.Background(), 123456)
	if err != nil {
		t.Fatalf("GameHighlights returned error: %v", err)
	}
	if !strings.Contains(rendered, "Earlier highlight (00:30)") {
		t.Fatalf("rendered output missing earlier highlight: %s", rendered)
	}
	if !strings.Contains(rendered, "https://video.example/earlier") {
		t.Fatalf("rendered output missing preferred playback URL: %s", rendered)
	}
	if !strings.Contains(rendered, "https://video.example/later") {
		t.Fatalf("rendered output missing fallback playback URL: %s", rendered)
	}
}
