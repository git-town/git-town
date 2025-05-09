Feature: create a new top-level feature branch in a clean workspace using the "compress" sync strategy when the branch has compressible commits

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE  |
      | feature | local, origin | commit 1 |
      | feature | local, origin | commit 2 |
    And the current branch is "feature"
    And Git setting "git-town.sync-feature-strategy" is "compress"
    When I run "git-town hack new"

  Scenario: result
    # TODO: shouln't it compress the commits here
    Then Git Town runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git fetch --prune --tags |
      |         | git checkout main        |
      | main    | git checkout -b new      |
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE  |
      | feature | local, origin | commit 1 |
      |         |               | commit 2 |
    And this lineage exists now
      | BRANCH  | PARENT |
      | feature | main   |
      | new     | main   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND              |
      | new     | git checkout feature |
      | feature | git branch -D new    |
    And the initial commits exist now
    And the initial branches and lineage exist now
