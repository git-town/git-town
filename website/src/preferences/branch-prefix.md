# Branch prefix

When set, Git Town automatically adds this prefix to branches it creates.

For example, with a branch prefix of `kevgo-`:

- `git town hack example` creates branch `kevgo-example`
- `git town append child` creates branch `kevgo-child`
- `git town prepend parent` creates branch `kevgo-parent`
- `git town rename other` renames the current branch to `kevgo-other`

If the branch name you provide already includes the configured prefix, Git Town
won't add it again. For instance, with prefix `kevgo-`, running
`git town hack kevgo-example` creates `kevgo-example` (not
`kevgo-kevgo-example`).

## configure in config file

In the [config file](../configuration-file.md), define the branch prefix within
the `[create]` section:

```toml
[create]
branch-prefix = "kevgo-"
```

## configure in Git metadata

To manually set the branch prefix, use the following command:

```wrap
git config [--global] git-town.branch-prefix 'kevgo-'
```

The optional `--global` flag applies this setting to all Git repositories on
your local machine. When not present, the setting applies to the current repo.

## environment variable

You can configure the branch prefix by setting the `GIT_TOWN_BRANCH_PREFIX`
environment variable.

If you want to use your GitHub username as the branch prefix, please set this
environment variable through something like this:

```bash
export GIT_TOWN_BRANCH_PREFIX=$(gh api user --jq '.login')
```
