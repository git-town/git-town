Feature: sync the current perennial branch using the ff-only sync strategy after the tracking branch was deleted at the remote

  Background:
    Given a Git repo with origin
    And the branches
      | NAME       | TYPE      | LOCATIONS     |
      | production | perennial | local, origin |
    And the commits
      | BRANCH     | LOCATION      | MESSAGE      |
      | production | local, origin | first commit |
    And the current branch is "production"
    And Git setting "git-town.sync-perennial-strategy" is "ff-only"
    And origin deletes the "production" branch
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH     | COMMAND                  |
      | production | git fetch --prune --tags |
      |            | git checkout main        |
      | main       | git branch -D production |
      |            | git push --tags          |
    And the current branch is now "main"
    And these branches exist now
      | REPOSITORY    | BRANCHES |
      | local, origin | main     |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                        |
      | main   | git branch production {{ sha 'first commit' }} |
      |        | git checkout production                        |
    And the current branch is still "production"
    And the initial commits exist now
    And the initial branches and lineage exist now
