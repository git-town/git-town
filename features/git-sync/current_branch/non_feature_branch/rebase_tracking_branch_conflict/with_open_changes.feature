Feature: git sync: resolving conflicts between the current non-feature branch and its tracking branch (with open changes)

  As a developer syncing a non-feature branch that conflicts with its tracking branch
  I want to be given the choice to resolve the conflicts or abort
  So that I can finish the operation as planned or postpone it to a better time.


  Background:
    Given I have branches named "qa" and "production"
    And my non-feature branches are configured as "qa" and "production"
    And I am on the "qa" branch
    And the following commits exist in my repository
      | BRANCH | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT               |
      | qa     | local    | conflicting local commit  | conflicting_file | local conflicting content  |
      |        | remote   | conflicting remote commit | conflicting_file | remote conflicting content |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git sync`


  Scenario: result
    Then it runs the Git commands
      | BRANCH | COMMAND              |
      | qa     | git fetch --prune    |
      |        | git stash -u         |
      |        | git rebase origin/qa |
    And I get the error
      """
      To abort, run "git sync --abort".
      To continue after you have resolved the conflicts, run "git sync --continue".
      To skip the sync of the 'qa' branch, run "git sync --skip".
      """
    And my repo has a rebase in progress
    And I don't have an uncommitted file with name: "uncommitted"


  Scenario: aborting
    When I run `git sync --abort`
    Then it runs the Git commands
      | BRANCH | COMMAND            |
      | qa     | git rebase --abort |
      |        | git stash pop      |
    And I am still on the "qa" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And there is no rebase in progress
    And I am left with my original commits


  Scenario: continuing without resolving the conflicts
    When I run `git sync --continue`
    Then it runs no Git commands
    And I get the error "You must resolve the conflicts before continuing the git sync"
    And my uncommitted file "uncommitted" is still stashed away
    And my repo still has a rebase in progress


  Scenario: continuing after resolving the conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `git sync --continue`
    Then it runs the Git commands
      | BRANCH | COMMAND               |
      | qa     | git rebase --continue |
      |        | git push              |
      |        | git push --tags       |
      |        | git stash pop         |
    And I am still on the "qa" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And now I have the following commits
      | BRANCH | LOCATION         | MESSAGE                   | FILE NAME        |
      | qa     | local and remote | conflicting remote commit | conflicting_file |
      |        |                  | conflicting local commit  | conflicting_file |
    And now I have the following committed files
      | BRANCH | NAME             | CONTENT          |
      | qa     | conflicting_file | resolved content |


  Scenario: continuing after resolving the conflicts and continuing the rebase
    Given I resolve the conflict in "conflicting_file"
    When I run `git rebase --continue; git sync --continue`
    Then it runs the Git commands
      | BRANCH | COMMAND         |
      | qa     | git push        |
      |        | git push --tags |
      |        | git stash pop   |
    And I am still on the "qa" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And now I have the following commits
      | BRANCH | LOCATION         | MESSAGE                   | FILE NAME        |
      | qa     | local and remote | conflicting remote commit | conflicting_file |
      |        |                  | conflicting local commit  | conflicting_file |
    And now I have the following committed files
      | BRANCH | NAME             | CONTENT          |
      | qa     | conflicting_file | resolved content |
