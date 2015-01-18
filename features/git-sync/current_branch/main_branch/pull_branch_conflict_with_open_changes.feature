Feature: git sync: resolving conflicting remote branch updates when syncing the main branch (with open changes)

  As a developer syncing the main branch when it conflicts with its tracking branch
  I want to be given the choice to resolve the conflicts or abort
  So that I can finish the operation as planned or postpone it to a better time.


  Background:
    Given I am on the "main" branch
    And the following commits exist in my repository
      | BRANCH | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT               |
      | main   | remote   | conflicting remote commit | conflicting_file | remote conflicting content |
      |        | local    | conflicting local commit  | conflicting_file | local conflicting content  |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    And I run `git sync`, it errors


  @finishes-with-non-empty-stash
  Scenario: result
    Then it runs the Git commands
      | BRANCH | COMMAND                |
      | main   | git fetch --prune      |
      | main   | git stash -u           |
      | main   | git rebase origin/main |
    And my repo has a rebase in progress
    And I don't have an uncommitted file with name: "uncommitted"


  Scenario: aborting
    When I run `git sync --abort`
    Then it runs the Git commands
      | BRANCH | COMMAND            |
      | HEAD   | git rebase --abort |
      | main   | git stash pop      |
    And I am still on the "main" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And there is no rebase in progress
    And I am left with my original commits


  @finishes-with-non-empty-stash
  Scenario: continuing without resolving conflicts
    When I run `git sync --continue`, it errors
    Then it runs no Git commands
    And I get the error "You must resolve the conflicts before continuing the git sync"
    And I don't have an uncommitted file with name: "uncommitted"
    And my repo still has a rebase in progress


  Scenario: continuing after resolving conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `git sync --continue`
    Then it runs the Git commands
      | BRANCH | COMMAND               |
      | HEAD   | git rebase --continue |
      | main   | git push              |
      | main   | git push --tags       |
      | main   | git stash pop         |
    And I am still on the "main" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And now I have the following commits
      | BRANCH | LOCATION         | MESSAGE                   | FILE NAME        |
      | main   | local and remote | conflicting remote commit | conflicting_file |
      |        |                  | conflicting local commit  | conflicting_file |
    And now I have the following committed files
      | BRANCH | FILES            | CONTENT          |
      | main   | conflicting_file | resolved content |


  Scenario: continuing after resolving conflicts and continuing the rebase
    Given I resolve the conflict in "conflicting_file"
    When I run `git rebase --continue; git sync --continue`
    Then it runs the Git commands
      | BRANCH | COMMAND         |
      | main   | git push        |
      | main   | git push --tags |
      | main   | git stash pop   |
    And I am still on the "main" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And now I have the following commits
      | BRANCH | LOCATION         | MESSAGE                   | FILE NAME        |
      | main   | local and remote | conflicting remote commit | conflicting_file |
      |        |                  | conflicting local commit  | conflicting_file |
    And now I have the following committed files
      | BRANCH | FILES            | CONTENT          |
      | main   | conflicting_file | resolved content |
