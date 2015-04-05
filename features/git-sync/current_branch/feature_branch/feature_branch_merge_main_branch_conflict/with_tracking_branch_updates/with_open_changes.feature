Feature: git sync: resolving conflicts between the current feature branch and the main branch (with tracking branch updates and open changes)

  Background:
    Given I have a feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE                    | FILE NAME        | FILE CONTENT    |
      | main    | local    | conflicting main commit    | conflicting_file | main content    |
      | feature | local    | conflicting feature commit | conflicting_file | feature content |
      |         | remote   | feature commit             | feature_file     | feature content |
    And I am on the "feature" branch
    And I have an uncommitted file
    When I run `git sync`


  @finishes-with-non-empty-stash
  Scenario: result
    Then it runs the Git commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune                  |
      |         | git stash -u                       |
      |         | git checkout main                  |
      | main    | git rebase origin/main             |
      |         | git push                           |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
    And I get the error
      """
      To abort, run "git sync --abort".
      To continue after you have resolved the conflicts, run "git sync --continue".
      To skip the sync of the 'feature' branch, run "git sync --skip".
      """
    And I am still on the "feature" branch
    And I don't have an uncommitted file with name: "uncommitted"
    And my repo has a merge in progress


  Scenario: aborting
    When I run `git sync --abort`
    Then it runs the Git commands
      | BRANCH  | COMMAND                                                  |
      | feature | git merge --abort                                        |
      |         | git reset --hard <%= sha 'conflicting feature commit' %> |
      |         | git checkout main                                        |
      | main    | git checkout feature                                     |
      | feature | git stash pop                                            |
    And I am still on the "feature" branch
    And I again have my uncommitted file
    And there is no merge in progress
    And I still have the following commits
      | BRANCH  | LOCATION         | MESSAGE                    | FILE NAME        | FILE CONTENT    |
      | main    | local and remote | conflicting main commit    | conflicting_file | main content    |
      | feature | local            | conflicting feature commit | conflicting_file | feature content |
      |         | remote           | feature commit             | feature_file     | feature content |


  @finishes-with-non-empty-stash
  Scenario: continuing without resolving the conflicts
    When I run `git sync --continue`
    Then it runs no Git commands
    And I get the error "You must resolve the conflicts before continuing the git sync"
    And I am still on the "feature" branch
    And I don't have an uncommitted file with name: "uncommitted"
    And my repo still has a merge in progress


  Scenario: continuing after resolving the conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `git sync --continue`
    Then it runs the Git commands
      | BRANCH  | COMMAND              |
      | feature | git commit --no-edit |
      |         | git push             |
      |         | git stash pop        |
    And I am still on the "feature" branch
    And I again have my uncommitted file
    And I still have the following commits
      | BRANCH  | LOCATION         | MESSAGE                                                    | FILE NAME        |
      | main    | local and remote | conflicting main commit                                    | conflicting_file |
      | feature | local and remote | conflicting feature commit                                 | conflicting_file |
      |         |                  | feature commit                                             | feature_file     |
      |         |                  | Merge remote-tracking branch 'origin/feature' into feature |                  |
      |         |                  | conflicting main commit                                    | conflicting_file |
      |         |                  | Merge branch 'main' into feature                           |                  |
    And I still have the following committed files
      | BRANCH  | NAME             | CONTENT          |
      | main    | conflicting_file | main content     |
      | feature | conflicting_file | resolved content |
      | feature | feature_file     | feature content  |


  Scenario: continuing after resolving the conflicts and comitting
    Given I resolve the conflict in "conflicting_file"
    When I run `git commit --no-edit; git sync --continue`
    Then it runs the Git commands
      | BRANCH  | COMMAND       |
      | feature | git push      |
      |         | git stash pop |
    And I am still on the "feature" branch
    And I again have my uncommitted file
    And I still have the following commits
      | BRANCH  | LOCATION         | MESSAGE                                                    | FILE NAME        |
      | main    | local and remote | conflicting main commit                                    | conflicting_file |
      | feature | local and remote | conflicting feature commit                                 | conflicting_file |
      |         |                  | feature commit                                             | feature_file     |
      |         |                  | Merge remote-tracking branch 'origin/feature' into feature |                  |
      |         |                  | conflicting main commit                                    | conflicting_file |
      |         |                  | Merge branch 'main' into feature                           |                  |
    And I still have the following committed files
      | BRANCH  | NAME             | CONTENT          |
      | main    | conflicting_file | main content     |
      | feature | conflicting_file | resolved content |
      | feature | feature_file     | feature content  |
