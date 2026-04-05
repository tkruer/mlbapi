---
layout: default
title: Getting started
---

# Getting started

## Requirements

- Go 1.22+

## Install

```bash
go get github.com/panthersinsights/mlbapi
```

## Basic usage

Create a client, then call endpoint methods:

```go
ctx := context.Background()
client := pkg.NewClient(nil)

schedule, err := client.Schedule(ctx, pkg.ScheduleRequest{
    SportID: []int{1},
})
if err != nil {
    // handle error
}

fmt.Println(schedule.TotalItems)
```

## Testing locally

```bash
make test
```

## Next steps

- Read the [API overview](./api)
- Review the [development guide](./development)
