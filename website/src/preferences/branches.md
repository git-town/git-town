# Branches

Git provides infrastructure to store many useful operations to create, merge, push/pull, and remove branches. Being an intentionally generic tool, it leaves how exactly you use Git branches and structure your development workflow completely up to you.

Git Town, a layer of additional commands above Git, provides predefined meaning and workflow automation for Git branches in typical usage scenarios.

At a high level, Git Town distinguishes long-lived from short-lived branches. Branches that live forever are called _perennial branches_. Typical names for perennial branches are `main`, `master`, `development`, `production`, or `staging`. A special type of perennial branch is the _main branch_. That's the branch from which short-lived branches are cut by default, and into which short-lived branches are merged by default.

Short-lived branches typically exist to develop code. They are cut from a perennial branch and merged back into the same perennial branch. Git Town distinguishes short-lived branches that you own vs those that you don't own.

Branches owned by you:
- feature branch: a branch that you currently work on. Git Town keeps it up to date for you.
- prototype branch: a feature branch at the early stage of development, before it is ready to be pushed to the development remote
- parked branch: a feature branch that you own but don't actively develop and aren't interested in syncing

Branches owned by somebody else:
- contribution branches: a feature branch owned by somebody else that you are contributing code to, but don't want to sync, ship, or delete
- observed branches: a feature branch owned by somebody else that you just are reviewing, but aren't contributing code to

Git Town offers powerful configuration settings to give each existing and new branch the correct type.

You tell Git Town which of your branches are perennial by providing a simple list of branch names, or a regular expression, and all branches whose name matches it are considered perennial.
The main branch is automatically marked perennial.
