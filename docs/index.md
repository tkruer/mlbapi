---
layout: default
title: mlbapi
---

# mlbapi

A Go client for MLB Stats API data including schedules, teams, game feeds, box scores, and stats.

## Why this project

`mlbapi` wraps frequently-used MLB Stats API endpoints with strongly typed Go models and tests so that applications can consume baseball data reliably.

## Highlights

- Idiomatic Go package (`pkg`) with typed request/response models.
- Unit tests and integration tests for core endpoint workflows.
- Generated endpoint constants for maintainable path usage.

## Quick start

```bash
go get github.com/panthersinsights/mlbapi
```

```go
package main

import (
    "context"
    "fmt"

    "github.com/panthersinsights/mlbapi/pkg"
)

func main() {
    ctx := context.Background()
    client := pkg.NewClient(nil)

    teams, err := client.Teams(ctx)
    if err != nil {
        panic(err)
    }

    fmt.Println("teams:", len(teams.Teams))
}
```

## Documentation

- [Getting started](./getting-started)
- [API overview](./api)
- [Development guide](./development)

## Repository

Source code and issue tracking live in the GitHub repository:

- <https://github.com/panthersinsights/mlbapi>
