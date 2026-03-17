# Propose headless

The _propose-headless_ preference controls whether the
[propose](../commands/propose.md) command opens a browser window after creating
a proposal. When enabled, proposals are created without opening a browser. This
is useful for CI/CD pipelines and headless environments.

## options

When set to `true`, proposals are created without opening a browser. When set to
`false` (default), `propose` opens a browser window after creating a proposal.

## configure in config file

In the [config file](../configuration-file.md):

```toml
[propose]
headless = true
```

## configure in Git metadata

In Git metadata:

```bash
git config [--global] git-town.propose-headless true
```

## configure via environment variable

```bash
export GIT_TOWN_PROPOSE_HEADLESS=true
```
