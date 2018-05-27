#### NAME

sync - updates the current branch with all relevant changes

#### SYNOPSIS

```
git town sync [--all]
```

#### DESCRIPTION

Synchronizes the current branch with the rest of the world.

When run on a feature branch

* syncs all ancestor branches
* pulls updates for the current branch
* merges the parent branch into the current branch
* pushes the current branch

When run on the main branch or a perennial branch

* pulls and pushes updates for the current branch
* pushes tags

Additionally, when there is a remote upstream,
the main branch is synced with its upstream counterpart.

#### OPTIONS

```
--all
    Syncs all local branches
```
