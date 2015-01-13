Feature: git sync: handling conflicting remote branch updates when syncing a non-feature branch (with open changes)

  As a developer syncing a non-feature branch that conflicts with its tracking branch
  I want to be given the choice to resolve the conflicts or abort
  So that I can finish the operation as planned or postpone it to a better time.


  Background:
    Given non-feature branch configuration "qa, production"
    And I am on the "qa" branch
    And the following commits exist in my repository
      | BRANCH | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT               |
      | qa     | remote   | conflicting remote commit | conflicting_file | remote conflicting content |
      |        | local    | conflicting local commit  | conflicting_file | local conflicting content  |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    And I run `git sync` while allowing errors


  @finishes-with-non-empty-stash
  Scenario: result
    Then it runs the Git commands
      | BRANCH | COMMAND              |
      | qa     | git fetch --prune    |
      | qa     | git stash -u         |
      | qa     | git rebase origin/qa |
    And my repo has a rebase in progress
    And I don't have an uncommitted file with name: "uncommitted"


  Scenario: aborting
    When I run `git sync --abort`
    Then it runs the Git commands
      | BRANCH | COMMAND            |
      | HEAD   | git rebase --abort |
      | qa     | git stash pop      |
    And I am still on the "qa" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And there is no rebase in progress
    And I am left with my original commits


  @finishes-with-non-empty-stash
  Scenario: continuing without resolving conflicts
    When I run `git sync --continue` while allowing errors
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
      | qa     | git push              |
      | qa     | git push --tags       |
      | qa     | git stash pop         |
    And I am still on the "qa" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And now I have the following commits
      | BRANCH | LOCATION         | MESSAGE                   | FILE NAME        |
      | qa     | local and remote | conflicting remote commit | conflicting_file |
      |        |                  | conflicting local commit  | conflicting_file |
    And now I have the following committed files
      | BRANCH | FILES            | CONTENT          |
      | qa     | conflicting_file | resolved content |


  Scenario: continuing after resolving conflicts and continuing the rebase
    Given I resolve the conflict in "conflicting_file"
    When I run `git rebase --continue; git sync --continue`
    Then it runs the Git commands
      | BRANCH | COMMAND         |
      | qa     | git push        |
      | qa     | git push --tags |
      | qa     | git stash pop   |
    And I am still on the "qa" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And now I have the following commits
      | BRANCH | LOCATION         | MESSAGE                   | FILE NAME        |
      | qa     | local and remote | conflicting remote commit | conflicting_file |
      |        |                  | conflicting local commit  | conflicting_file |
    And now I have the following committed files
      | BRANCH | FILES            | CONTENT          |
      | qa     | conflicting_file | resolved content |
