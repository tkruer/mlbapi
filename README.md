# mlbapi

`mlbapi` is a Go client for the public MLB Stats API at `https://statsapi.mlb.com/api/`.

It provides two layers:

- Typed helper functions for common workflows such as schedules, standings, lookups, box scores, highlights, roster output, and player summaries.
- Raw endpoint access through a reusable `Client`, generic `Get`, and endpoint-specific `Call...` methods.

The Go package lives in [`/mlbapi/pkg`](/mlbapi/pkg).

## Installation

```bash
go get github.com/tkruer/mlbapi
```

Import it as:

```go
import mlbapi "github.com/tkruer/mlbapi/pkg"
```

Current package version:

- `v0.0.1`

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"

	mlbapi "github.com/tkruer/mlbapi/pkg"
)

func main() {
	games, err := mlbapi.Schedule(context.Background(), mlbapi.ScheduleOptions{
		Date: "2026-04-03",
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, game := range games {
		fmt.Println(game.Summary)
	}
}
```

## Typed Helpers

The package exposes higher-level helpers for the most common MLB API use cases.

Schedule and game data:

- `Schedule`
- `LastGame`
- `NextGame`
- `Linescore`
- `Boxscore`
- `BoxscoreDataForGame`
- `GameScoringPlays`
- `GameScoringPlayData`
- `GameHighlights`
- `GameHighlightData`

Player and team lookups:

- `LookupPlayer`
- `LookupTeam`
- `PlayerStatData`
- `PlayerStats`
- `Roster`

League and metadata helpers:

- `StandingsData`
- `Standings`
- `TeamLeaderData`
- `TeamLeaders`
- `LeagueLeaderData`
- `LeagueLeaders`
- `Meta`
- `Notes`

These helpers return typed values where that improves usability. Examples include:

- `Schedule` returns `[]ScheduleGame`
- `LookupPlayer` returns `[]PlayerLookup`
- `LookupTeam` returns `[]Team`
- `StandingsData` returns `[]DivisionStandings`

## Team Constants and Metadata

The package includes a stable `TeamID` type plus constants for all current MLB clubs.

Examples:

- `mlbapi.TeamLosAngelesDodgers`
- `mlbapi.TeamNewYorkYankees`
- `mlbapi.TeamBostonRedSox`
- `mlbapi.TeamAthletics`

You can resolve team metadata without making an API call:

```go
dodgers, ok := mlbapi.TeamByID(mlbapi.TeamLosAngelesDodgers)
if !ok {
	log.Fatal("team not found")
}

fmt.Println(dodgers.Name)
fmt.Println(dodgers.Abbreviation)
fmt.Println(dodgers.TeamCode)
```

Available helpers:

- `AllTeams()`
- `TeamByID(id TeamID)`
- `TeamByAbbreviation(abbreviation string)`
- `TeamByName(name string)`

## Raw Endpoint Access

If you want full control, use the generic client methods directly.

```go
client := mlbapi.NewClient()

payload, err := client.Get(context.Background(), "teams", mlbapi.Params{
	"sportIds": 1,
})
if err != nil {
	log.Fatal(err)
}

fmt.Println(len(payload["teams"].([]any)))
```

You can also call endpoint-specific methods:

```go
payload, err := client.CallSchedule(context.Background(), mlbapi.Params{
	"sportId": 1,
	"date":    "2026-04-03",
})
if err != nil {
	log.Fatal(err)
}

_ = payload
```

For advanced cases, `GetForce` will include otherwise-invalid query parameters instead of dropping them.

## Client Configuration

Use `NewClient` with options when you need a custom base URL or HTTP client.

```go
httpClient := &http.Client{
	Timeout: 15 * time.Second,
}

client := mlbapi.NewClient(
	mlbapi.WithHTTPClient(httpClient),
)
```

Available options:

- `WithBaseURL`
- `WithHTTPClient`

## Notes and Meta

`Notes(endpoint)` returns local documentation about supported path/query parameters and any endpoint-specific notes.

```go
notes, err := mlbapi.Notes("schedule")
if err != nil {
	log.Fatal(err)
}

fmt.Println(notes)
```

`Meta(type)` returns metadata values from the Stats API for supported meta categories such as:

- `positions`
- `pitchTypes`
- `statGroups`
- `standingsTypes`

## Testing and Development

Local development commands:

- `make fmt`
- `make fmt-check`
- `make vet`
- `make lint`
- `make test`
- `make test-race`
- `make integration-test`
- `make cover`
- `make ci`

Default tests are hermetic and do not call the live MLB API.

To run the live integration tests that verify real `200 OK` responses from `statsapi.mlb.com`:

```bash
make integration-test
```

The live integration workflow is also available as a manual GitHub Actions run in [`/mlbapi/.github/workflows/integration.yml`](/mlbapi/.github/workflows/integration.yml).

## Releases and pkg.go.dev

Go module releases are fully automated on `main` and use conventional commits to compute semantic version bumps.

- Branch CI runs automatically on pushes and pull requests for `develop` and `main`
- A push to `main` triggers the release workflow in [`/mlbapi/.github/workflows/release.yml`](/mlbapi/.github/workflows/release.yml)
- The release workflow verifies formatting, linting, tests, race tests, and coverage before creating a release
- It finds commits since the latest `v*.*.*` tag and applies semver rules:
  - `BREAKING CHANGE:` or `type!:` (`feat!:` / `fix!:`) → major bump
  - `feat:` → minor bump
  - `fix:` → patch bump
  - other commit types do not create a release
- When a releasable commit exists, the workflow creates and pushes the next tag, creates a GitHub Release with generated notes, and requests module indexing on `proxy.golang.org`

`mlbapi.Version` is set to `"dev"`; Git tags are the canonical source of released versions.

## Package Layout

Main package files:

- [pkg/client.go](/mlbapi/pkg/client.go)
- [pkg/types.go](/mlbapi/pkg/types.go)
- [pkg/schedule.go](/mlbapi/pkg/schedule.go)
- [pkg/boxscore.go](/mlbapi/pkg/boxscore.go)
- [pkg/stats.go](/mlbapi/pkg/stats.go)
- [pkg/teams.go](/mlbapi/pkg/teams.go)

## Status

This package currently targets the public MLB Stats API and is organized for direct application use:

- typed helper surface for common tasks
- raw endpoint access when you need the full API
- local tooling for formatting, linting, testing, and coverage
- opt-in live integration tests for real API verification

## Project documentation site

This repository includes a GitHub Pages documentation site under `docs/`.

- Source docs: `docs/`
- Deployment workflow: `.github/workflows/docs.yml`

After merging to `main`, the docs site is automatically deployed via GitHub Actions.
