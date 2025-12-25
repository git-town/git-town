# Typical development workflow

The following four Git Town commands automate the typical development workflow:

- You start hacking by running [git town hack](./commands/hack.md) to create a
  feature branch.
- While coding you run [git town sync](./commands/sync.md) to keep your feature
  branch up to date with commits that you or other developers make into the main
  branch. This prevents your feature branch from deviating too much from the
  main code line.
- If your team does pull requests, you can run
  [git town propose](./commands/propose.md) to create a new pull request.
- [git town ship](./commands/ship.md) delivers the feature branch.
