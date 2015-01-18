Feature: git sync: resolving conflicting remote feature branch updates when syncing a feature branch with open changes

  As a developer syncing a feature branch that conflicts with the tracking branch
  I want to be given the choice to resolve the conflicts or abort
  So that I can finish the operation as planned or postpone it to a better time.


  Background:
    Given I have a feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT               |
      | feature | remote   | remote conflicting commit | conflicting_file | remote conflicting content |
      |         | local    | local conflicting commit  | conflicting_file | local conflicting content  |
    And I am on the "feature" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    And I run `git sync`, it errors


  @finishes-with-non-empty-stash
  Scenario: result
    Then it runs the Git commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune                  |
      | feature | git stash -u                       |
      | feature | git checkout main                  |
      | main    | git rebase origin/main             |
      | main    | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
    And I am still on the "feature" branch
    And I don't have an uncommitted file with name: "uncommitted"
    And my repo has a merge in progress


  Scenario: aborting
    When I run `git sync --abort`
    Then it runs the Git commands
      | BRANCH  | COMMAND              |
      | feature | git merge --abort    |
      | feature | git checkout main    |
      | main    | git checkout feature |
      | feature | git stash pop        |
    And I am still on the "feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And there is no merge in progress
    And I am left with my original commits


  @finishes-with-non-empty-stash
  Scenario: continuing without resolving conflicts
    When I run `git sync --continue`, it errors
    Then I get the error "You must resolve the conflicts before continuing the git sync"
    And I am still on the "feature" branch
    And I don't have an uncommitted file with name: "uncommitted"
    And my repo still has a merge in progress


  Scenario: continuing after resolving conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `git sync --continue`
    Then it runs the Git commands
      | BRANCH  | COMMAND                  |
      | feature | git commit --no-edit     |
      | feature | git merge --no-edit main |
      | feature | git push                 |
      | feature | git stash pop            |
    And I am still on the "feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And now I have the following commits
      | BRANCH  | LOCATION         | MESSAGE                                                    | FILE NAME        |
      | feature | local and remote | local conflicting commit                                   | conflicting_file |
      |         |                  | remote conflicting commit                                  | conflicting_file |
      |         |                  | Merge remote-tracking branch 'origin/feature' into feature |                  |
    And now I have the following committed files
      | BRANCH  | FILES            | CONTENT          |
      | feature | conflicting_file | resolved content |


  Scenario: continuing after resolving conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `git commit --no-edit; git sync --continue`
    Then it runs the Git commands
      | BRANCH  | COMMAND                  |
      | feature | git merge --no-edit main |
      | feature | git push                 |
      | feature | git stash pop            |
    And I am still on the "feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And now I have the following commits
      | BRANCH  | LOCATION         | MESSAGE                                                    | FILE NAME        |
      | feature | local and remote | local conflicting commit                                   | conflicting_file |
      |         |                  | remote conflicting commit                                  | conflicting_file |
      |         |                  | Merge remote-tracking branch 'origin/feature' into feature |                  |
    And now I have the following committed files
      | BRANCH  | FILES            | CONTENT          |
      | feature | conflicting_file | resolved content |
