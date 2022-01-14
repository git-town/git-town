# Typical development workflow

The following four Git Town commands automate the typical development workflow:

- You start hacking by running [git hack](./commands/hack.md) to create a
  feature branch.
- While coding and committing you run [git sync](./commands/sync.md) to keep
  your feature branch up to date with commits that other developers make into
  the main branch. This prevents your feature branch from deviating too much
  from the main code line.
- If your team does pull requests, you can run
  [git new-pull-request](./commands/new-pull-request.md) to create a new pull
  request.
- You run [git ship](./commands/ship.md) to deliver the feature branch.
