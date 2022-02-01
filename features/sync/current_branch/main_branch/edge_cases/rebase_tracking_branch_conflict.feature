Feature: handle conflicts between the main branch and its tracking branch when syncing the main branch

  Background:
    Given I am on the "main" branch
    And the following commits exist in my repo
      | BRANCH | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT               |
      | main   | local    | conflicting local commit  | conflicting_file | local conflicting content  |
      |        | remote   | conflicting remote commit | conflicting_file | remote conflicting content |
    And my workspace has an uncommitted file
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
    And my uncommitted file is stashed

  Scenario: aborting
    When I run "git-town abort"
    Then it runs the commands
      | BRANCH | COMMAND            |
      | main   | git rebase --abort |
      |        | git stash pop      |
    And I am still on the "main" branch
    And my workspace still contains my uncommitted file
    And there is no rebase in progress
    And my repo is left with my original commits

  Scenario: continuing without resolving the conflicts
    When I run "git-town continue"
    Then it runs no commands
    And it prints the error:
      """
      you must resolve the conflicts before continuing
      """
    And my uncommitted file is stashed
    And my repo still has a rebase in progress

  Scenario: continuing after resolving the conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run "git-town continue" and close the editor
    Then it runs the commands
      | BRANCH | COMMAND               |
      | main   | git rebase --continue |
      |        | git push              |
      |        | git push --tags       |
      |        | git stash pop         |
    And I am still on the "main" branch
    And my workspace still contains my uncommitted file
    And all branches are now synchronized
    And my repo now has the following committed files
      | BRANCH | NAME             | CONTENT          |
      | main   | conflicting_file | resolved content |

  Scenario: continuing after resolving the conflicts and continuing the rebase
    Given I resolve the conflict in "conflicting_file"
    And I run "git rebase --continue" and close the editor
    When I run "git-town continue"
    Then it runs the commands
      | BRANCH | COMMAND         |
      | main   | git push        |
      |        | git push --tags |
      |        | git stash pop   |
    And I am still on the "main" branch
    And my workspace still contains my uncommitted file
    And all branches are now synchronized
    And my repo now has the following committed files
      | BRANCH | NAME             | CONTENT          |
      | main   | conflicting_file | resolved content |
