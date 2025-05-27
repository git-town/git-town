# Branches

At a high level, Git Town distinguishes long-lived from short-lived Git branches. Branches that live forever are called _perennial branches_. Typical names for perennial branches are `main`, `master`, `development`, `production`, or `staging`. A special type of perennial branch is the _main branch_. That's the branch from which short-lived branches are cut by default, and into which short-lived branches are merged by default.

Short-lived branches typically exist to develop code. They are cut from a perennial branch and merged back into the same perennial branch. Git Town distinguishes short-lived branches that you own vs those that you don't own.

Branches owned by you:
- **feature branch:** a branch that you currently work on. Git Town keeps it up to date for you
- **prototype branch:** a feature branch at an early stage of development, before it is ready to be pushed to the development remote
- **parked branch:** a feature branch that you own but don't actively work on, and aren't interested in syncing

Branches owned by somebody else:
- **contribution branches:** a feature branch owned by somebody else that you are contributing code to, but don't want to sync, ship, or delete
- **observed branches:** a feature branch owned by somebody else that you are reviewing, but aren't contributing code to

Git Town offers powerful configuration settings to give each existing and new branch the correct type.

- perennial-branches: a list of branch names that are perennial
- perennial-regex: all branches matching this regular expression are considered perennial
- main-branch: The main branch is automatically considered perennial.
- contribution-regex: all branches matching this regular expression are considered contribution branches
- observed-regex: all branches matching this regular expression are considered observed branches
- new-branch-type: the type that new branches created by git town hack, append, or prepend should have

- git town contribute: marks a branch as a contribution branch
- git town observe: marks a branch as an observed branch
- git town park: marks a branch as a parked branch
- git town prototype: create a new prototype branch or mark an existing branch as a prototype branch
- git town hack: create a new feature branch or mark an existing branch as a feature branch
