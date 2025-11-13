# Lint Cached Connectors

This linter verifies that cached connector implementations in `internal/forge` implement the same interfaces as their uncached counterparts.

## Purpose

The codebase has pairs of connectors:
- **Uncached connectors**: Direct API/CLI connectors (e.g., `APIConnector`, `Connector`)
- **Cached connectors**: Wrappers that add caching (e.g., `CachedAPIConnector`, `CachedConnector`)

Each cached connector should implement exactly the same interfaces as its uncached counterpart to ensure they are interchangeable.

## How It Works

The linter:
1. Parses Go source files for connector pairs
2. Extracts interface implementations from type check comments (e.g., `var _ forgedomain.ProposalFinder = apiConnector`)
3. Compares interfaces between cached and uncached versions
4. Reports any missing interface implementations

## Usage

```bash
# Build and run the linter
make test

# Build only
make build

# Clean build artifacts
make clean
```

## Connector Pairs Checked

- `bitbucketcloud`: `APIConnector` / `CachedAPIConnector`
- `bitbucketdatacenter`: `APIConnector` / `CachedAPIConnector`
- `forgejo`: `APIConnector` / `CachedAPIConnector`
- `gitea`: `APIConnector` / `CachedAPIConnector`
- `github`: `APIConnector` / `CachedAPIConnector`
- `gitlab`: `APIConnector` / `CachedAPIConnector`
- `gh`: `Connector` / `CachedConnector`
- `glab`: `Connector` / `CachedConnector`

## Adding New Connector Pairs

To add a new connector pair to the linter, edit `main.go` and add an entry to the `connectorPairs` slice:

```go
{
    Package:      "myforge",
    UncachedFile: "internal/forge/myforge/api_connector.go",
    UncachedType: "APIConnector",
    CachedFile:   "internal/forge/myforge/cached_api_connector.go",
    CachedType:   "CachedAPIConnector",
},
```
