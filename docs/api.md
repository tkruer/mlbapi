---
layout: default
title: API overview
---

# API overview

The package exposes a typed client in `pkg` for major MLB Stats API domains.

## Client

- `pkg.NewClient(httpClient *http.Client) *Client`

## Endpoint groups

- Teams
- Schedule
- Game / box score
- Player and team stats

## Key package files

- `pkg/client.go` for transport and request logic.
- `pkg/types.go` for shared response/request models.
- `pkg/*` domain files for endpoint-specific methods.

## Tips

- Reuse one `Client` across your app to share connection pooling.
- Pass `context.Context` through each call for cancellation/timeouts.
- Use integration tests when validating endpoint compatibility.
