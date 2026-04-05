---
layout: default
title: Development guide
---

# Development guide

## Build and test

```bash
make test
```

## Code generation

Some endpoint constants are generated. Regenerate and format before committing:

```bash
go generate ./...
go test ./...
```

## CI and release

The repository includes GitHub Actions workflows for CI and release automation.

## Contributing checklist

1. Add/update tests for behavior changes.
2. Run `go test ./...` locally.
3. Update docs in `docs/` or `README.md` when APIs change.
