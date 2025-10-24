# Order

This setting allows you to change how Git Town orders branches it displays.

Allowed values:

- **asc** sort branches in natural sort order, ascending (default)
- **desc** sort branches in natural sort order, descending

## CLI flag

You can override this setting per command using:

- `--order=asc` to force ascending order
- `--order=desc` to force descending order

## config file

```toml
[branches]
order = "<asc|desc>"
```

## Git metadata

To enable ordering branches in Git, run this command:

```wrap
git config [--global] git-town.order <asc|desc>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, this setting applies to the current Git repo.

## environment variable

You can configure branches ordering by setting the `GIT_TOWN_ORDER` environment
variable.
