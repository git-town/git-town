Feature: git sync: resolving conflicts between the main branch and its tracking branch when syncing the current feature branch (without open changes)

  (see ./with_open_changes.feature)


  Background:
    Given I have a feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT               |
      | main   | remote   | conflicting remote commit | conflicting_file | remote conflicting content |
      |        | local    | conflicting local commit  | conflicting_file | local conflicting content  |
    And I am on the "feature" branch
    When I run `git sync`


  Scenario: result
    Then it runs the Git commands
      | BRANCH  | COMMAND                |
      | feature | git fetch --prune      |
      | feature | git checkout main      |
      | main    | git rebase origin/main |
    And I get the error
      """
      To abort, run "git sync --abort".
      To continue after you have resolved the conflicts, run "git sync --continue".
      """
    And my repo has a rebase in progress


  Scenario: aborting
    When I run `git sync --abort`
    Then it runs the Git commands
      | BRANCH | COMMAND              |
      | main   | git rebase --abort   |
      | main   | git checkout feature |
    And I am still on the "feature" branch
    And there is no rebase in progress
    And I am left with my original commits


  Scenario: continuing without resolving conflicts
    When I run `git sync --continue`
    Then I get the error "You must resolve the conflicts before continuing the git sync"
    And my repo still has a rebase in progress
    And I don't have an uncommitted file with name: "uncommitted"


  Scenario: continuing after resolving conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `git sync --continue`
    Then it runs the Git commands
      | BRANCH  | COMMAND                            |
      | main    | git rebase --continue              |
      | main    | git push                           |
      | main    | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      | feature | git merge --no-edit main           |
      | feature | git push                           |
    And I am still on the "feature" branch
    And now I have the following commits
      | BRANCH  | LOCATION         | MESSAGE                   | FILE NAME        |
      | main    | local and remote | conflicting remote commit | conflicting_file |
      |         |                  | conflicting local commit  | conflicting_file |
      | feature | local and remote | conflicting remote commit | conflicting_file |
      |         |                  | conflicting local commit  | conflicting_file |
    And now I have the following committed files
      | BRANCH  | FILES            | CONTENT          |
      | main    | conflicting_file | resolved content |
      | feature | conflicting_file | resolved content |


  Scenario: continuing after resolving conflicts and continuing the rebase
    Given I resolve the conflict in "conflicting_file"
    When I run `git rebase --continue; git sync --continue`
    Then it runs the Git commands
      | BRANCH  | COMMAND                            |
      | main    | git push                           |
      | main    | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      | feature | git merge --no-edit main           |
      | feature | git push                           |
    And I am still on the "feature" branch
    And now I have the following commits
      | BRANCH  | LOCATION         | MESSAGE                   | FILE NAME        |
      | main    | local and remote | conflicting remote commit | conflicting_file |
      |         |                  | conflicting local commit  | conflicting_file |
      | feature | local and remote | conflicting remote commit | conflicting_file |
      |         |                  | conflicting local commit  | conflicting_file |
    And now I have the following committed files
      | BRANCH  | FILES            | CONTENT          |
      | main    | conflicting_file | resolved content |
      | feature | conflicting_file | resolved content |
