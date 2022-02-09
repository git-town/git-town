Feature: handle conflicts between the main branch and its tracking branch when syncing the main branch

  Background:
    Given the current branch is "main"
    And the commits
      | BRANCH | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT   |
      | main   | local    | conflicting local commit  | conflicting_file | local content  |
      |        | origin   | conflicting origin commit | conflicting_file | origin content |
    And an uncommitted file
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
      |        | git add -A               |
      |        | git stash                |
      |        | git rebase origin/main   |
    And it prints the error:
      """
      To abort, run "git-town abort".
      To continue after having resolved conflicts, run "git-town continue".
      """
    And my repo now has a rebase in progress
    And the uncommitted file is stashed

  Scenario: abort
    When I run "git-town abort"
    Then it runs the commands
      | BRANCH | COMMAND            |
      | main   | git rebase --abort |
      |        | git stash pop      |
    And the current branch is still "main"
    And the uncommitted file still exists
    And there is no rebase in progress anymore
    And now the initial commits exist

  Scenario: continue with unresolved conflict
    When I run "git-town continue"
    Then it runs no commands
    And it prints the error:
      """
      you must resolve the conflicts before continuing
      """
    And the uncommitted file is stashed
    And my repo still has a rebase in progress

  Scenario: resolve and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue" and close the editor
    Then it runs the commands
      | BRANCH | COMMAND               |
      | main   | git rebase --continue |
      |        | git push              |
      |        | git push --tags       |
      |        | git stash pop         |
    And all branches are now synchronized
    And the current branch is still "main"
    And the uncommitted file still exists
    And my repo now has these committed files
      | BRANCH | NAME             | CONTENT          |
      | main   | conflicting_file | resolved content |

  Scenario: resolve, finish the rebase, and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git rebase --continue" and close the editor
    And I run "git-town continue"
    Then it runs the commands
      | BRANCH | COMMAND         |
      | main   | git push        |
      |        | git push --tags |
      |        | git stash pop   |
    And all branches are now synchronized
    And the current branch is still "main"
    And the uncommitted file still exists
    And my repo now has these committed files
      | BRANCH | NAME             | CONTENT          |
      | main   | conflicting_file | resolved content |
