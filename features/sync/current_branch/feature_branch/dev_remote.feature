Feature: sync a branch to a custom dev-remote

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS |
      | branch | feature | main   | local     |
    And the commits
      | BRANCH | LOCATION | MESSAGE |
      | branch | local    | commit  |
    And I rename the "origin" remote to "fork"
    And Git setting "git-town.dev-remote" is "fork"
    And the current branch is "branch"
    And I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | branch | git fetch --prune --tags |
      |        | git push -u fork branch  |
    And the branches are now
      | REPOSITORY  | BRANCHES     |
      | local, fork | main, branch |
    And all branches are now synchronized
    And these commits exist now
      | BRANCH | LOCATION    | MESSAGE |
      | branch | local, fork | commit  |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND               |
      | branch | git push fork :branch |
    And the branches are now
      | REPOSITORY | BRANCHES     |
      | local      | main, branch |
      | fork       | main         |
    And the initial commits exist now
