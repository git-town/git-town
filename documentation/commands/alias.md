#### NAME

alias - add or remove the default git aliases


#### SYNOPSIS

```
git town alias [true | false]
```

#### DESCRIPTION

Adding the default git aliases removes the need for `town` for the following commands (append, hack, kill, new-pull-request, prepend, prune-branches, rename-branch, repo, ship, and sync). Example: `git append` becomes equivalent to `git town append`.

When adding aliases, no existing aliases will be overwritten.
