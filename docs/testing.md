# Testing

Run unit tests with coverage:

```bash
go test ./... -coverprofile /tmp/peky.cover
go tool cover -func /tmp/peky.cover | tail -n 1
```

Race tests:

```bash
go test ./... -race
```

CLI smoke run (builds `./bin/peky`, starts daemon, runs core commands):

```bash
scripts/cli-smoke.sh
```

CLI stress battery (builds `./bin/peky`, isolated runtime/config, heavy concurrency):

```bash
scripts/cli-stress.sh
```

Run tool loops (Codex/Claude) during stress:

```bash
RUN_TOOLS=1 scripts/cli-stress.sh
```

Performance + profiling tools live in `docs/performance.md`.

GitHub Actions runs gofmt checks, go vet, go test with coverage, and race on Linux.

Nightly self-hosted stress (Geekom A8 / Linux runner):

- Workflow: `.github/workflows/nightly-stress.yml`
- Runs `scripts/cli-stress.sh` on `self-hosted, Linux, X64`

## Release Safety Checklist (pre-release, no tagging)

Run these before a release candidate to avoid regressions:

1) **Clean build & tests**

```bash
go test ./...
go test ./... -race
```

2) **Coverage sanity check**

```bash
go test ./... -coverprofile /tmp/peky.cover
go tool cover -func /tmp/peky.cover | tail -n 1
```

3) **Perf baseline (dev-only)**

See `docs/performance.md` for `scripts/perf-profiler` and `scripts/perf-bench`.

4) **Update perf summary**

```bash
# Write/refresh: testdata/performance-tests/perf-summary-YYYYMMDD.md
```

5) **Verify dev-only flags are gated**

- Profiling hooks must be behind build tags or env gates.
- No profiler endpoints enabled by default.
