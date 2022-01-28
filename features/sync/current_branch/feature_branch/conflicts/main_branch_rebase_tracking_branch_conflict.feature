Feature: resolving conflicts between the main branch and its tracking branch when syncing the current feature branch

  Background:
    Given my repo has a feature branch named "feature"
    And the following commits exist in my repo
      | BRANCH | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT               |
      | main   | local    | conflicting local commit  | conflicting_file | local conflicting content  |
      |        | remote   | conflicting remote commit | conflicting_file | remote conflicting content |
    And I am on the "feature" branch
    And my workspace has an uncommitted file
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git fetch --prune --tags |
      |         | git add -A               |
      |         | git stash                |
      |         | git checkout main        |
      | main    | git rebase origin/main   |
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
      | BRANCH  | COMMAND              |
      | main    | git rebase --abort   |
      |         | git checkout feature |
      | feature | git stash pop        |
    And I am still on the "feature" branch
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
    And my repo still has a rebase in progress
    And my uncommitted file is stashed

  Scenario: continuing after resolving the conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run "git-town continue" and close the editor
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | main    | git rebase --continue              |
      |         | git push                           |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
      |         | git push                           |
      |         | git stash pop                      |
    And I am still on the "feature" branch
    And my workspace still contains my uncommitted file
    And my repo now has the following commits
      | BRANCH  | LOCATION      | MESSAGE                   | FILE NAME        |
      | main    | local, remote | conflicting remote commit | conflicting_file |
      |         |               | conflicting local commit  | conflicting_file |
      | feature | local, remote | conflicting remote commit | conflicting_file |
      |         |               | conflicting local commit  | conflicting_file |
    And my repo now has the following committed files
      | BRANCH  | NAME             | CONTENT          |
      | main    | conflicting_file | resolved content |
      | feature | conflicting_file | resolved content |

  Scenario: continuing after resolving the conflicts and continuing the rebase
    Given I resolve the conflict in "conflicting_file"
    And I run "git rebase --continue" and close the editor
    When I run "git-town continue"
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | main    | git push                           |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
      |         | git push                           |
      |         | git stash pop                      |
    And I am still on the "feature" branch
    And my workspace still contains my uncommitted file
    And my repo now has the following commits
      | BRANCH  | LOCATION      | MESSAGE                   | FILE NAME        |
      | main    | local, remote | conflicting remote commit | conflicting_file |
      |         |               | conflicting local commit  | conflicting_file |
      | feature | local, remote | conflicting remote commit | conflicting_file |
      |         |               | conflicting local commit  | conflicting_file |
    And my repo now has the following committed files
      | BRANCH  | NAME             | CONTENT          |
      | main    | conflicting_file | resolved content |
      | feature | conflicting_file | resolved content |
