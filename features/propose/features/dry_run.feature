@skipWindows
Feature: dry-run proposing changes

  Background: proposing changes
    Given a Git repo with origin
    And the branch
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And tool "open" is installed
    And the current branch is "feature"
    And the origin is "git@github.com:git-town/git-town.git"

  @this
  Scenario: a PR for this branch exists already
    Given a proposal for this branch exists at "https://github.com/git-town/git-town/pull/123"
    When I run "git-town propose --dry-run"
    Then it runs the commands
      | BRANCH  | COMMAND                                            |
      | feature | git fetch --prune --tags                           |
      | <none>  | looking for proposal online ... ok                 |
      | feature | git checkout main                                  |
      | main    | git rebase origin/main                             |
      |         | git checkout feature                               |
      | feature | git merge --no-edit --ff origin/feature            |
      |         | git merge --no-edit --ff main                      |
      | <none>  | open https://github.com/git-town/git-town/pull/123 |
    And the current branch is still "feature"
    And the initial branches and lineage exist

  Scenario: there is no PR for this branch yet
