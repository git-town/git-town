Feature: handle conflicts between the current feature branch and its tracking branch

  Background:
    Given my repo has a feature branch "feature"
    And my repo contains the commits
      | BRANCH  | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT               |
      | feature | local    | local conflicting commit  | conflicting_file | local conflicting content  |
      |         | remote   | remote conflicting commit | conflicting_file | remote conflicting content |
    And I am on the "feature" branch
    And my workspace has an uncommitted file
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune --tags           |
      |         | git add -A                         |
      |         | git stash                          |
      |         | git checkout main                  |
      | main    | git rebase origin/main             |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
    And it prints the error:
      """
      To abort, run "git-town abort".
      To continue after having resolved conflicts, run "git-town continue".
      To continue by skipping the current branch, run "git-town skip".
      """
    And I am still on the "feature" branch
    And my uncommitted file is stashed
    And my repo now has a merge in progress

  Scenario: aborting
    When I run "git-town abort"
    Then it runs the commands
      | BRANCH  | COMMAND              |
      | feature | git merge --abort    |
      |         | git checkout main    |
      | main    | git checkout feature |
      | feature | git stash pop        |
    And I am still on the "feature" branch
    And my workspace still contains my uncommitted file
    And there is no merge in progress
    And my repo is left with my original commits

  Scenario: continuing without resolving the conflicts
    When I run "git-town continue"
    Then it runs no commands
    And it prints the error:
      """
      you must resolve the conflicts before continuing
      """
    And I am still on the "feature" branch
    And my uncommitted file is stashed
    And my repo still has a merge in progress

  Scenario: continuing after resolving the conflicts
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue"
    Then it runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git commit --no-edit     |
      |         | git merge --no-edit main |
      |         | git push                 |
      |         | git stash pop            |
    And I am still on the "feature" branch
    And my workspace still contains my uncommitted file
    And all branches are now synchronized
    And my repo now has the following committed files
      | BRANCH  | NAME             | CONTENT          |
      | feature | conflicting_file | resolved content |

  Scenario: continuing after resolving the conflicts and committing
    When I resolve the conflict in "conflicting_file"
    And I run "git commit --no-edit"
    And I run "git-town continue"
    Then it runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git merge --no-edit main |
      |         | git push                 |
      |         | git stash pop            |
    And I am still on the "feature" branch
    And my workspace still contains my uncommitted file
    And all branches are now synchronized
    And my repo now has the following committed files
      | BRANCH  | NAME             | CONTENT          |
      | feature | conflicting_file | resolved content |
