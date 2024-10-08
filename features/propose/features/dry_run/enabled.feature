@skipWindows
Feature: dry-run proposing changes

  Background: proposing changes
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And tool "open" is installed
    And the current branch is "feature"
    And the origin is "git@github.com:git-town/git-town.git"

  Scenario: a PR for this branch exists already
    Given a proposal for this branch exists at "https://github.com/git-town/git-town/pull/123"
    When I run "git-town propose --dry-run"
    Then it runs the commands
      | BRANCH  | COMMAND                                            |
      | feature | git fetch --prune --tags                           |
      | <none>  | Looking for proposal online ... ok                 |
      |         | open https://github.com/git-town/git-town/pull/123 |
    And the current branch is still "feature"
    And the initial branches and lineage exist now

  Scenario: there is no PR for this branch yet
    Given a proposal for this branch does not exist
    When I run "git-town propose --dry-run"
    Then it runs the commands
      | BRANCH  | COMMAND                                                            |
      | feature | git fetch --prune --tags                                           |
      | <none>  | Looking for proposal online ... ok                                 |
      | feature | git checkout main                                                  |
      | main    | git rebase origin/main --no-update-refs                            |
      |         | git checkout feature                                               |
      | feature | git merge --no-edit --ff origin/feature                            |
      |         | git merge --no-edit --ff main                                      |
      | <none>  | open https://github.com/git-town/git-town/compare/feature?expand=1 |
    And the current branch is still "feature"
    And the initial branches and lineage exist now
