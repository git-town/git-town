# Typical development workflow

The following four Git Town commands automate the typical development workflow:

- You start hacking by running [git hack](./commands/hack.md) to create a
  feature branch.
- While coding you run [git sync](./commands/sync.md) to keep your feature
  branch up to date with commits that you or other developers make into the main
  branch. This prevents your feature branch from deviating too much from the
  main code line.
- If your team does pull requests, you can run
  [git new-pull-request](./commands/new-pull-request.md) to create a new pull
  request.
- [git ship](./commands/ship.md) delivers the feature branch.
