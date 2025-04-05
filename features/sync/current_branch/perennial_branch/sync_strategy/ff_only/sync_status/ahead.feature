Feature: sync the current perennial branch using the ff-only sync strategy when ahead of the tracking branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME       | TYPE      | LOCATIONS     |
      | production | perennial | local, origin |
    And the commits
      | BRANCH     | LOCATION | MESSAGE      |
      | production | local    | first commit |
    And the current branch is "production"
    And Git setting "git-town.sync-perennial-strategy" is "ff-only"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH     | COMMAND                  |
      | production | git fetch --prune --tags |
    And Git Town prints the error:
      """
      cannot sync branch "production" because it has unpushed local commits
      """
    And the initial branches and lineage exist now
    And the initial commits exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the initial commits exist now
    And the initial branches and lineage exist now
