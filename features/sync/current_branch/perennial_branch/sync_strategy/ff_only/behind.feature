Feature: sync the current perennial branch using the rebase sync strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME       | TYPE      | LOCATIONS     |
      | production | perennial | local, origin |
    And the commits
      | BRANCH     | LOCATION | MESSAGE      |
      | production | origin   | first commit |
    And the current branch is "production"
    And Git setting "git-town.sync-perennial-strategy" is "ff-only"
    And inspect the repo
    When I run "git-town sync"

  @debug
  @this
  Scenario: result
    Then Git Town runs the commands
      | BRANCH     | COMMAND                  |
      | production | git fetch --prune --tags |
    And the current branch is still "production"
    And the initial branches and lineage exist now
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE       |
      | main   | local, origin | main commit   |
      | qa     | local, origin | origin commit |
      |        |               | local commit  |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the current branch is still "qa"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE       |
      | main   | local, origin | main commit   |
      | qa     | local, origin | origin commit |
      |        |               | local commit  |
    And the initial branches and lineage exist now
