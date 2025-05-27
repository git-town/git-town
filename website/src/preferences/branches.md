# Branches

At a high level, Git Town distinguishes long-lived from short-lived Git
branches.

## Long-lived branches

Branches that live forever are called _perennial branches_. Typical names for
perennial branches are `main`, `master`, `development`, `production`, or
`staging`. Amongst these, the _main branch_ holds a special role: it is the the
default base from which short-lived branches are cut, and into which short-lived
branches are merged.

## Short-lived branches

Short-lived branches typically used for active development. They are typically
created from a perennial branch and merged back into the same perennial branch.
They can also form a hierarchy of branches called a
[stack](../stacked-changes.md). Git Town distinguishes short-lived branches that
you own vs those that you don't own.

### Short-lived branches owned by you

- **feature branch:** a branch that you do work on, Git Town keeps it up to date
  for you
- **prototype branch:** an early-stage feature branch, not ready to be pushed to
  a shared remote
- **parked branch:** a feature branch you own but aren't actively working on,
  Git Town doesn't sync it to reduce noise

### Short-lived branches owned by others

- **contribution branches:** somebody else's feature branch that you are
  contributing code to, but no lifecycle events like sync, ship, or delete
- **observed branches:** somebody else's feature branch that you review but
  aren't contributing code to

## Configuring branch types

Git Town offers powerful configuration settings to automatically classify local
branches:

- [branches.main](main-branch.md): automatically treated as perennial
- [branches.perennials](perennial-branches.md) explicit list of perennial
  branches
- [branches.perennial-regex](perennial-regex.md) all branches matching this
  regular expression are considered perennial
- [branches.contribution-regex](contribution-regex.md): all branches matching
  this regular expression are considered contribution branches
- [branches.observed-regex](observed-regex.md): all branches matching this
  regular expression are considered observed branches
- [create.new-branch-type](new-branch-type.md) defines the type that branches
  you create via commands like [git town hack](../commands/hack.md),
  [append](../commands/append.md), or [prepend](../commands/prepend.md)

## Manually setting branch types

You can override the branch type for each branch using one of these commands:

- [git town contribute](../commands/contribute.md): mark a branch as a
  contribution branch
- [git town observe](../commands/observe.md): mark a branch as observed
- [git town park](../commands/park.md): mark a branch as parked
- [git town prototype](../commands/prototype.md): create or convert a branch to
  a prototype branch
- [git town hack](../commands/hack.md): create or convert a branch to a feature
  branch
