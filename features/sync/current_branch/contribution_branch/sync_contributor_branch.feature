Feature: sync the current contribution branch

  Background:
    Given the current branch is a contribution branch "contribution"
    And the commits
      | BRANCH       | LOCATION      | MESSAGE       | FILE NAME   |
      | main         | local, origin | main commit   | main_file   |
      | contribution | local         | local commit  | local_file  |
      |              | origin        | origin commit | origin_file |
    When I run "git-town sync"

  @this
  Scenario: result
    Then it runs the commands
      | BRANCH       | COMMAND                        |
      | contribution | git fetch --prune --tags       |
      |              | git rebase origin/contribution |
      |              | git push                       |
      |              | git push --tags                |
    And the current branch is still "contribution"
    And these commits exist now
      | BRANCH       | LOCATION      | MESSAGE       |
      | main         | local, origin | main commit   |
      | contribution | local, origin | origin commit |
      |              | local, origin | local commit  |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH       | COMMAND                                              |
      | contribution | git reset --hard {{ sha-before-run 'local commit' }} |
    And the current branch is still "contribution"
    And the initial commits exist
    And the initial branches and lineage exist
