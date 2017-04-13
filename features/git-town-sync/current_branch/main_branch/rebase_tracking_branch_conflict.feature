Feature: gt sync: resolving conflicts between the main branch and its tracking branch when syncing the main branch

  As a developer syncing the main branch when it conflicts with its tracking branch
  I want to be given the choice to resolve the conflicts or abort
  So that I can finish the operation as planned or postpone it to a better time.


  Background:
    Given I am on the "main" branch
    And the following commits exist in my repository
      | BRANCH | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT               |
      | main   | local    | conflicting local commit  | conflicting_file | local conflicting content  |
      |        | remote   | conflicting remote commit | conflicting_file | remote conflicting content |
    And I have an uncommitted file
    When I run `gt sync`


  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                |
      | main   | git fetch --prune      |
      |        | git add -A             |
      |        | git stash              |
      |        | git rebase origin/main |
    And I get the error
      """
      To abort, run "gt sync --abort".
      To continue after you have resolved the conflicts, run "gt sync --continue".
      """
    And my repo has a rebase in progress
    And my uncommitted file is stashed


  Scenario: aborting
    When I run `gt sync --abort`
    Then it runs the commands
      | BRANCH | COMMAND            |
      | main   | git rebase --abort |
      |        | git stash pop      |
    And I am still on the "main" branch
    And I still have my uncommitted file
    And there is no rebase in progress
    And I am left with my original commits


  Scenario: continuing without resolving the conflicts
    When I run `gt sync --continue`
    Then it runs no commands
    And I get the error "you must resolve the conflicts before continuing"
    And my uncommitted file is stashed
    And my repo still has a rebase in progress


  Scenario: continuing after resolving the conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `gt sync --continue`
    Then it runs the commands
      | BRANCH | COMMAND               |
      | main   | git rebase --continue |
      |        | git push              |
      |        | git push --tags       |
      |        | git stash pop         |
    And I am still on the "main" branch
    And I still have my uncommitted file
    And now I have the following commits
      | BRANCH | LOCATION         | MESSAGE                   | FILE NAME        |
      | main   | local and remote | conflicting remote commit | conflicting_file |
      |        |                  | conflicting local commit  | conflicting_file |
    And now I have the following committed files
      | BRANCH | NAME             | CONTENT          |
      | main   | conflicting_file | resolved content |


  Scenario: continuing after resolving the conflicts and continuing the rebase
    Given I resolve the conflict in "conflicting_file"
    When I run `git rebase --continue; gt sync --continue`
    Then it runs the commands
      | BRANCH | COMMAND         |
      | main   | git push        |
      |        | git push --tags |
      |        | git stash pop   |
    And I am still on the "main" branch
    And I still have my uncommitted file
    And now I have the following commits
      | BRANCH | LOCATION         | MESSAGE                   | FILE NAME        |
      | main   | local and remote | conflicting remote commit | conflicting_file |
      |        |                  | conflicting local commit  | conflicting_file |
    And now I have the following committed files
      | BRANCH | NAME             | CONTENT          |
      | main   | conflicting_file | resolved content |
