Feature: git-town sync: resolving conflicts between the current perennial branch and its tracking branch

  As a developer syncing a perennial branch that conflicts with its tracking branch
  I want to be given the choice to resolve the conflicts or abort
  So that I can finish the operation as planned or postpone it to a better time.


  Background:
    Given my repository has the perennial branches "production" and "qa"
    And I am on the "qa" branch
    And the following commits exist in my repository
      | BRANCH | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT               |
      | qa     | local    | conflicting local commit  | conflicting_file | local conflicting content  |
      |        | remote   | conflicting remote commit | conflicting_file | remote conflicting content |
    And my workspace has an uncommitted file
    When I run `git-town sync`


  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | qa     | git fetch --prune --tags |
      |        | git add -A               |
      |        | git stash                |
      |        | git rebase origin/qa     |
    And it prints the error:
      """
      To abort, run "git-town abort".
      To continue after having resolved conflicts, run "git-town continue".
      To continue by skipping the current branch, run "git-town skip".
      """
    And my repo has a rebase in progress
    And my uncommitted file is stashed


  Scenario: aborting
    When I run `git-town abort`
    Then it runs the commands
      | BRANCH | COMMAND            |
      | qa     | git rebase --abort |
      |        | git stash pop      |
    And I am still on the "qa" branch
    And my workspace still contains my uncommitted file
    And there is no rebase in progress
    And my repository is left with my original commits


  Scenario: continuing without resolving the conflicts
    When I run `git-town continue`
    Then it runs no commands
    And it prints the error "You must resolve the conflicts before continuing"
    And my uncommitted file is stashed
    And my repo still has a rebase in progress


  Scenario: continuing after resolving the conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `git-town continue`
    Then it runs the commands
      | BRANCH | COMMAND               |
      | qa     | git rebase --continue |
      |        | git push              |
      |        | git push --tags       |
      |        | git stash pop         |
    And I am still on the "qa" branch
    And my workspace still contains my uncommitted file
    And now my repository has the following commits
      | BRANCH | LOCATION         | MESSAGE                   | FILE NAME        |
      | qa     | local and remote | conflicting remote commit | conflicting_file |
      |        |                  | conflicting local commit  | conflicting_file |
    And now my repository has the following committed files
      | BRANCH | NAME             | CONTENT          |
      | qa     | conflicting_file | resolved content |


  Scenario: continuing after resolving the conflicts and continuing the rebase
    Given I resolve the conflict in "conflicting_file"
    When I run `git rebase --continue; git-town continue`
    Then it runs the commands
      | BRANCH | COMMAND         |
      | qa     | git push        |
      |        | git push --tags |
      |        | git stash pop   |
    And I am still on the "qa" branch
    And my workspace still contains my uncommitted file
    And now my repository has the following commits
      | BRANCH | LOCATION         | MESSAGE                   | FILE NAME        |
      | qa     | local and remote | conflicting remote commit | conflicting_file |
      |        |                  | conflicting local commit  | conflicting_file |
    And now my repository has the following committed files
      | BRANCH | NAME             | CONTENT          |
      | qa     | conflicting_file | resolved content |
