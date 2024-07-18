Feature: inactive parked branches don't get synced

  Background:
    Given a Git repo clone
    And the branches
      | NAME   | TYPE   | PARENT | LOCATIONS     |
      | parked | parked | main   | local, origin |
    And the commits
      | BRANCH | LOCATION | MESSAGE              |
      | main   | local    | local main commit    |
      |        | origin   | origin main commit   |
      | parked | local    | local parked commit  |
      |        | origin   | origin parked commit |
    And the current branch is "main"
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
      |        | git rebase origin/main   |
      |        | git push                 |
      |        | git push --tags          |
    And the current branch is still "main"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE              |
      | main   | local, origin | origin main commit   |
      |        |               | local main commit    |
      | parked | local         | local parked commit  |
      |        | origin        | origin parked commit |

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "main"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE              |
      | main   | local, origin | origin main commit   |
      |        |               | local main commit    |
      | parked | local         | local parked commit  |
      |        | origin        | origin parked commit |
    And the initial branches and lineage exist
