Feature: sync the current perennial branch using the ff-only sync strategy when already in sync with the tracking branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME       | TYPE      | LOCATIONS     |
      | production | perennial | local, origin |
    And the commits
      | BRANCH     | LOCATION      | MESSAGE      |
      | production | local, origin | first commit |
    And Git setting "git-town.sync-perennial-strategy" is "ff-only"
    And the current branch is "production"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH     | COMMAND                  |
      | production | git fetch --prune --tags |
      |            | git push --tags          |
    And the initial branches and lineage exist now
    And the initial commits exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the initial branches and lineage exist now
    And the initial commits exist now
