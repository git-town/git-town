Feature: git sync: resolving conflicts between the current feature branch and its tracking branch (without open changes)

  (see ./pull_feature_branch_conflict_with_open_changes.feature)


  Background:
    Given I have a feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT               |
      | feature | local    | local conflicting commit  | conflicting_file | local conflicting content  |
      |         | remote   | remote conflicting commit | conflicting_file | remote conflicting content |
    And I am on the "feature" branch
    When I run `git sync`


  Scenario: result
    Then it runs the Git commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune                  |
      | feature | git checkout main                  |
      | main    | git rebase origin/main             |
      | main    | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
    And I get the error
      """
      To abort, run "git sync --abort".
      To continue after you have resolved the conflicts, run "git sync --continue".
      To skip the sync of the 'feature' branch, run "git sync --skip".
      """
    And I am still on the "feature" branch
    And my repo has a merge in progress


  Scenario: aborting
    When I run `git sync --abort`
    Then it runs the Git commands
      | BRANCH  | COMMAND              |
      | feature | git merge --abort    |
      | feature | git checkout main    |
      | main    | git checkout feature |
    And I am still on the "feature" branch
    And there is no merge in progress
    And I am left with my original commits


  Scenario: continuing without resolving the conflicts
    When I run `git sync --continue`
    Then I get the error "You must resolve the conflicts before continuing the git sync"
    And I am still on the "feature" branch
    And my repo still has a merge in progress


  Scenario: continuing after resolving the conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `git sync --continue`
    Then it runs the Git commands
      | BRANCH  | COMMAND                  |
      | feature | git commit --no-edit     |
      | feature | git merge --no-edit main |
      | feature | git push                 |
    And I am still on the "feature" branch
    And now I have the following commits
      | BRANCH  | LOCATION         | MESSAGE                                                    | FILE NAME        |
      | feature | local and remote | local conflicting commit                                   | conflicting_file |
      |         |                  | remote conflicting commit                                  | conflicting_file |
      |         |                  | Merge remote-tracking branch 'origin/feature' into feature |                  |
    And now I have the following committed files
      | BRANCH  | NAME             | CONTENT          |
      | feature | conflicting_file | resolved content |


  Scenario: continuing after resolving the conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `git commit --no-edit; git sync --continue`
    Then it runs the Git commands
      | BRANCH  | COMMAND                  |
      | feature | git merge --no-edit main |
      | feature | git push                 |
    And I am still on the "feature" branch
    And now I have the following commits
      | BRANCH  | LOCATION         | MESSAGE                                                    | FILE NAME        |
      | feature | local and remote | local conflicting commit                                   | conflicting_file |
      |         |                  | remote conflicting commit                                  | conflicting_file |
      |         |                  | Merge remote-tracking branch 'origin/feature' into feature |                  |
    And now I have the following committed files
      | BRANCH  | NAME             | CONTENT          |
      | feature | conflicting_file | resolved content |
