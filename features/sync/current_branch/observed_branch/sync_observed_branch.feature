Feature: sync the current observed branch

  Background:
    Given the current branch is an observed branch "observed"
    And the commits
      | BRANCH   | LOCATION      | MESSAGE       | FILE NAME   |
      | main     | local, origin | main commit   | main_file   |
      | observed | local         | local commit  | local_file  |
      |          | origin        | origin commit | origin_file |
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH   | COMMAND                    |
      | observed | git fetch --prune --tags   |
      |          | git rebase origin/observed |
      |          | git push --tags            |
    And the current branch is still "observed"
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE       |
      | main     | local, origin | main commit   |
      | observed | local, origin | origin commit |
      |          | local         | local commit  |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH   | COMMAND                                              |
      | observed | git reset --hard {{ sha-before-run 'local commit' }} |
    And the current branch is still "observed"
    And the initial commits exist
    And the initial branches and lineage exist
