Feature: sync the current parked branch with a tracking branch using the "merge" sync-feature strategy

  Background:
    Given a parked branch "parked"
    And the commits
      | BRANCH | LOCATION | MESSAGE              |
      | main   | local    | local main commit    |
      |        | origin   | origin main commit   |
      | parked | local    | local parked commit  |
      |        | origin   | origin parked commit |
    And the current branch is "main"
    When I run "git-town sync"

  @this
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
    Then it runs the commands
      | BRANCH | COMMAND                                                                    |
      | parked | git reset --hard {{ sha 'local parked commit' }}                           |
      |        | git push --force-with-lease origin {{ sha 'origin parked commit' }}:parked |
    And the current branch is still "parked"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE              |
      | main   | local, origin | origin main commit   |
      |        |               | local main commit    |
      | parked | local         | local parked commit  |
      |        | origin        | origin parked commit |
    And the initial branches and lineage exist
