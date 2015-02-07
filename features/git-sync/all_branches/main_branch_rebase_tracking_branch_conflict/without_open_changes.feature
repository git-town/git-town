Feature: git sync --all: handling rebase conflicts between main branch and its tracking branch without open changes

  Background:
    Given I have a feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE            | FILE NAME        | FILE CONTENT        |
      | main    | local    | main local commit  | conflicting_file | main local content  |
      | main    | remote   | main remote commit | conflicting_file | main remote content |
      | feature | local    | feature commit     | feature_file     | feature content     |
    And I am on the "main" branch
    When I run `git sync --all`


  Scenario: result
    Then it runs the Git commands
      | BRANCH | COMMAND                |
      | main   | git fetch --prune      |
      | main   | git rebase origin/main |
    And I get the error
      """
      To abort, run "git sync --abort".
      To continue after you have resolved the conflicts, run "git sync --continue".
      """
    And my repo has a rebase in progress


  Scenario: aborting
    When I run `git sync --abort`
    Then it runs the Git commands
      | BRANCH | COMMAND            |
      | main   | git rebase --abort |
    And I end up on the "main" branch
    And I have the following commits
      | BRANCH  | LOCATION | MESSAGE            | FILE NAME        |
      | main    | local    | main local commit  | conflicting_file |
      | main    | remote   | main remote commit | conflicting_file |
      | feature | local    | feature commit     | feature_file     |


  Scenario: continuing without resolving conflicts
    When I run `git sync --continue`
    Then it runs no Git commands
    And I get the error "You must resolve the conflicts before continuing the git sync"
    And my repo still has a rebase in progress


  Scenario: continuing after resolving conflicts
    Given I resolve the conflict in "conflicting_file"
    And I run `git sync --continue`
    Then it runs the Git commands
      | BRANCH  | COMMAND                            |
      | main    | git rebase --continue              |
      | main    | git push                           |
      | main    | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      | feature | git merge --no-edit main           |
      | feature | git push                           |
      | feature | git checkout main                  |
    And I end up on the "main" branch
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
    And I end up on the "main" branch
    And I have the following commits
      | BRANCH  | LOCATION         | MESSAGE                          | FILE NAME        |
      | main    | local and remote | main remote commit               | conflicting_file |
      |         |                  | main local commit                | conflicting_file |
      | feature | local and remote | feature commit                   | feature_file     |
      |         |                  | main remote commit               | conflicting_file |
      |         |                  | main local commit                | conflicting_file |
      |         |                  | Merge branch 'main' into feature |                  |
