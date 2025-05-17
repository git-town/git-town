# Branches

Git provides infrastructure to store many useful operations to create, merge, push/pull, and remove branches. Being an intentionally generic tool, it leaves how exactly you use Git branches and structure your development workflow completely up to you.

Git Town, a layer of additional commands above Git, provides predefined meaning and workflow automation for Git branches in typical usage scenarios.

At a high level, Git Town distinguishes long-lived from short-lived branches. Branches that live forever are called _perennial_.



 while feature branches exist only for a short time, typically to develop code.

Perennial branches:
- main branch
- perennial branch

Git Town distinguishes feature branches that you own vs those that you don't own.

Owned branches:
- feature branches
- prototype branches: early-stage feature branches
- parked branches

Not owned branches:
- observed branches
- contribution branches
