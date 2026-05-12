# Release Process

This document is the single source of truth for releasing peky.

## Preconditions

- You are on `main` with a clean working tree.
- You have push access to the GitHub repo.
- GitHub releases and Homebrew tap updates are done by GitHub Actions.

## Required Tests (must pass)

Run these before any release (matches CI):

```bash
unformatted=$(git ls-files -z '*.go' | xargs -0 gofmt -l)
if [ -n "$unformatted" ]; then
  echo "gofmt needed on:"
  echo "$unformatted"
  exit 1
fi

go vet ./...
go test ./... -coverprofile=cover.out
go tool cover -func=cover.out | tail -n 1
go test ./... -race
```

## Release Steps

1) Pick the next version

2) Update docs:
- Move `CHANGELOG.md` “Unreleased” → the new version with today’s date.
- Ensure `README.md` matches current behavior.

3) Commit and tag:

```bash
git add -A
git commit -m "release: vX.Y.Z"
git tag vX.Y.Z
git push origin main --tags
```

4) GitHub Actions publishes the release:

- Tag push triggers the `release` workflow, which runs GoReleaser to create the GitHub release + upload assets, then updates the Homebrew tap formula.

Monitor runs (recommended):

```bash
gh run list -w release -L 3
gh run view --log --web
```

## Cut Next Release

Use this when you need to cut the next release and ensure Homebrew publishing completes.

1) Confirm the `release` workflow succeeded (Homebrew tap update is part of it).
2) If it fails, fix and re-run the exact workflow for the same tag:

```bash
gh run list -w release -L 3
gh run view --log --web
```

3) Verify install (recommended):

```bash
brew install --build-from-source regenrek/tap/peky
brew test regenrek/tap/peky
```

## Release Helper (recommended)

You can use the scripted helper to run the local parts (tests, tag, push) and trigger the GitHub Actions release pipeline:

```bash
scripts/release X.Y.Z
```

Dry run (no tag/push/release, tests still run):

```bash
scripts/release X.Y.Z --dry-run
```

The helper requires a clean working tree and push access to the repo.

## Post-Release Verification

```bash
brew info regenrek/tap/peky
```

## Notes

- The `release` workflow builds binaries into `dist/` and creates the GitHub release.
- The `release` workflow updates the Homebrew tap formula (`regenrek/homebrew-tap`) via `scripts/update-homebrew-tap`.
