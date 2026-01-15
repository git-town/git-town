@skipWindows
Feature: dry-run proposing changes

  Background:
    Given a Git repo with origin
    And the origin is "git@github.com:git-town/git-town.git"
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the current branch is "feature"
    And tool "open" is installed

  Scenario: a PR for this branch exists already
    Given the proposals
      | ID | SOURCE BRANCH | TARGET BRANCH | URL                                           |
      |  1 | feature       | main          | https://github.com/git-town/git-town/pull/123 |
    When I run "git-town propose --dry-run"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                            |
      | feature | git fetch --prune --tags                           |
      |         | Finding proposal from feature into main ... ok     |
      |         | open https://github.com/git-town/git-town/pull/123 |
    And the initial branches and lineage exist now

  Scenario: there is no PR for this branch yet
    When I run "git-town propose --dry-run"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                            |
      | feature | git fetch --prune --tags                                           |
      |         | Finding proposal from feature into main ... ok                     |
      |         | open https://github.com/git-town/git-town/compare/feature?expand=1 |
    And the initial branches and lineage exist now
