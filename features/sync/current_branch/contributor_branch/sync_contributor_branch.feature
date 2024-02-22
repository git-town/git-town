Feature: sync the current contributor branch

  Background:
    Given the current branch is a contributor branch "contributor"
    And the commits
      | BRANCH      | LOCATION      | MESSAGE       | FILE NAME   |
      | main        | local, origin | main commit   | main_file   |
      | contributor | local         | local commit  | local_file  |
      |             | origin        | origin commit | origin_file |
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH      | COMMAND                       |
      | contributor | git fetch --prune --tags      |
      |             | git rebase origin/contributor |
      |             | git push                      |
      |             | git push --tags               |
    And the current branch is still "contributor"
    And these commits exist now
      | BRANCH      | LOCATION      | MESSAGE       |
      | main        | local, origin | main commit   |
      | contributor | local, origin | origin commit |
      |             | local, origin | local commit  |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH      | COMMAND                                              |
      | contributor | git reset --hard {{ sha-before-run 'local commit' }} |
    And the current branch is still "contributor"
    And the initial commits exist
    And the initial branches and lineage exist
