# Branches

At a high level, Git Town distinguishes long-lived from short-lived Git branches. Branches that live forever are called _perennial branches_. Typical names for perennial branches are `main`, `master`, `development`, `production`, or `staging`. A special type of perennial branch is the _main branch_. That's the branch from which short-lived branches are cut by default, and into which short-lived branches are merged by default.

Short-lived branches typically exist to develop code. They are cut from a perennial branch and merged back into the same perennial branch. Git Town distinguishes short-lived branches that you own vs those that you don't own.

Branches owned by you:
- **feature branch:** a branch that you currently work on. Git Town keeps it up to date for you
- **prototype branch:** a feature branch at an early stage of development, before it is ready to be pushed to the development remote
- **parked branch:** a feature branch that you own but don't actively work on, and aren't interested in syncing

Branches owned by somebody else:
- **contribution branches:** a feature branch owned by somebody else, you are contributing code but no lifecycle events like sync, ship, or delete
- **observed branches:** a feature branch owned by somebody else, you are reviewing, but aren't contributing code

Git Town offers powerful configuration settings to give each existing and new branch the correct type.

- [main-branch](main-branch.md): The main branch is automatically considered perennial.
- [perennial-branches](perennial-branches.md) a list of branch names that are perennial
- [perennial-regex](perennial-regex.md) all branches matching this regular expression are considered perennial
- [contribution-regex](contribution-regex.md): all branches matching this regular expression are considered contribution branches
- [observed-regex](observed-regex.md): all branches matching this regular expression are considered observed branches
- [new-branch-type](new-branch-type.md): the type that branches you create through [git town hack](../commands/hack.md), [append](../commands/append.md), [prepend](../commands/prepend.md)


You can also override the branch type for each branch using these commands:

- [git town contribute](../commands/contribute.md) marks a branch as a contribution branch
- [git town observe](../commands/observe.md) marks a branch as an observed branch
- [git town park](../commands/park.md) marks a branch as a parked branch
- [git town prototype](../commands/prototype.md) creates a new prototype branch or mark an existing branch as a prototype branch
- [git town hack](../commands/hack.md) creates a new feature branch or mark an existing branch as a feature branch
