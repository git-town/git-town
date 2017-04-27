Feature: git-town sync: resolving conflicts between the main branch and its tracking branch when syncing the current feature branch

  As a developer syncing a feature branch when there are conflicts between the local and remote main branches
  I want to be given the choice to resolve the conflicts or abort
  So that I can finish the operation as planned or postpone it to a better time.


  Background:
    Given I have a feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT               |
      | main   | local    | conflicting local commit  | conflicting_file | local conflicting content  |
      |        | remote   | conflicting remote commit | conflicting_file | remote conflicting content |
    And I am on the "feature" branch
    And I have an uncommitted file
    When I run `git-town sync`


  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                |
      | feature | git fetch --prune      |
      |         | git add -A             |
      |         | git stash              |
      |         | git checkout main      |
      | main    | git rebase origin/main |
    And I get the error
      """
      To abort, run "git-town sync --abort".
      To continue after you have resolved the conflicts, run "git-town sync --continue".
      """
    And my repo has a rebase in progress
    And my uncommitted file is stashed


  Scenario: aborting
    When I run `git-town sync --abort`
    Then it runs the commands
      | BRANCH  | COMMAND              |
      | main    | git rebase --abort   |
      |         | git checkout feature |
      | feature | git stash pop        |
    And I am still on the "feature" branch
    And I still have my uncommitted file
    And there is no rebase in progress
    And I am left with my original commits


  Scenario: continuing without resolving the conflicts
    When I run `git-town sync --continue`
    Then I get the error "You must resolve the conflicts before continuing"
    And my repo still has a rebase in progress
    And my uncommitted file is stashed


  Scenario: continuing after resolving the conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `git-town sync --continue`
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
    And I still have my uncommitted file
    And now I have the following commits
      | BRANCH  | LOCATION         | MESSAGE                   | FILE NAME        |
      | main    | local and remote | conflicting remote commit | conflicting_file |
      |         |                  | conflicting local commit  | conflicting_file |
      | feature | local and remote | conflicting remote commit | conflicting_file |
      |         |                  | conflicting local commit  | conflicting_file |
    And now I have the following committed files
      | BRANCH  | NAME             | CONTENT          |
      | main    | conflicting_file | resolved content |
      | feature | conflicting_file | resolved content |


  Scenario: continuing after resolving the conflicts and continuing the rebase
    Given I resolve the conflict in "conflicting_file"
    When I run `git rebase --continue; git-town sync --continue`
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | main    | git push                           |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
      |         | git push                           |
      |         | git stash pop                      |
    And I am still on the "feature" branch
    And I still have my uncommitted file
    And now I have the following commits
      | BRANCH  | LOCATION         | MESSAGE                   | FILE NAME        |
      | main    | local and remote | conflicting remote commit | conflicting_file |
      |         |                  | conflicting local commit  | conflicting_file |
      | feature | local and remote | conflicting remote commit | conflicting_file |
      |         |                  | conflicting local commit  | conflicting_file |
    And now I have the following committed files
      | BRANCH  | NAME             | CONTENT          |
      | main    | conflicting_file | resolved content |
      | feature | conflicting_file | resolved content |
