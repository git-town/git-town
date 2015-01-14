Feature: git sync --all: handling rebase conflicts between main branch and its tracking branch with open changes

  Background:
    Given I have a feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE            | FILE NAME        | FILE CONTENT        |
      | main    | local    | main local commit  | conflicting_file | main local content  |
      | main    | remote   | main remote commit | conflicting_file | main remote content |
      | feature | local    | feature commit     | feature_file     | feature content     |
    And I am on the "main" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git sync --all` while allowing errors


  @finishes-with-non-empty-stash
  Scenario: result
    Then it runs the Git commands
      | BRANCH | COMMAND                |
      | main   | git fetch --prune      |
      | main   | git stash -u           |
      | main   | git rebase origin/main |
    And I don't have an uncommitted file with name: "uncommitted"
    And my repo has a rebase in progress


  Scenario: aborting
    When I run `git sync --abort`
    Then it runs the Git commands
      | BRANCH | COMMAND            |
      | HEAD   | git rebase --abort |
      | main   | git stash pop      |
    And I end up on the "main" branch
    And I again have an uncommitted file with name: "uncommitted" and content: "stuff"
    And I have the following commits
      | BRANCH  | LOCATION | MESSAGE            | FILE NAME        |
      | main    | local    | main local commit  | conflicting_file |
      | main    | remote   | main remote commit | conflicting_file |
      | feature | local    | feature commit     | feature_file     |


  @finishes-with-non-empty-stash
  Scenario: continuing without resolving conflicts
    When I run `git sync --continue` while allowing errors
    Then it runs no Git commands
    And I get the error "You must resolve the conflicts before continuing the git sync"
    And I don't have an uncommitted file with name: "uncommitted"
    And my repo still has a rebase in progress


  Scenario: continuing after resolving conflicts
    Given I resolve the conflict in "conflicting_file"
    And I run `git sync --continue`
    Then it runs the Git commands
      | BRANCH  | COMMAND                            |
      | HEAD    | git rebase --continue              |
      | main    | git push                           |
      | main    | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      | feature | git merge --no-edit main           |
      | feature | git push                           |
      | feature | git checkout main                  |
      | main    | git stash pop                      |
    And I end up on the "main" branch
    And I again have an uncommitted file with name: "uncommitted" and content: "stuff"
    And I have the following commits
      | BRANCH  | LOCATION         | MESSAGE                          | FILE NAME        |
      | main    | local and remote | main remote commit               | conflicting_file |
      |         |                  | main local commit                | conflicting_file |
      | feature | local and remote | feature commit                   | feature_file     |
      |         |                  | main remote commit               | conflicting_file |
      |         |                  | main local commit                | conflicting_file |
      |         |                  | Merge branch 'main' into feature |                  |


  Scenario: continuing after resolving conflicts and continuing the rebase
    Given I resolve the conflict in "conflicting_file"
    And I run `git rebase --continue; git sync --continue`
    Then it runs the Git commands
      | BRANCH  | COMMAND                            |
      | main    | git push                           |
      | main    | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      | feature | git merge --no-edit main           |
      | feature | git push                           |
      | feature | git checkout main                  |
      | main    | git stash pop                      |
    And I end up on the "main" branch
    And I again have an uncommitted file with name: "uncommitted" and content: "stuff"
    And I have the following commits
      | BRANCH  | LOCATION         | MESSAGE                          | FILE NAME        |
      | main    | local and remote | main remote commit               | conflicting_file |
      |         |                  | main local commit                | conflicting_file |
      | feature | local and remote | feature commit                   | feature_file     |
      |         |                  | main remote commit               | conflicting_file |
      |         |                  | main local commit                | conflicting_file |
      |         |                  | Merge branch 'main' into feature |                  |
