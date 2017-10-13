Feature: git-town sync: resolving conflicts between the current feature branch and the main branch (with tracking branch updates)

  Background:
    Given my repository has a feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE                    | FILE NAME        | FILE CONTENT    |
      | main    | local    | conflicting main commit    | conflicting_file | main content    |
      | feature | local    | conflicting feature commit | conflicting_file | feature content |
      |         | remote   | feature commit             | feature_file     | feature content |
    And I am on the "feature" branch
    And my workspace has an uncommitted file
    When I run `git-town sync`


  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune                  |
      |         | git add -A                         |
      |         | git stash                          |
      |         | git checkout main                  |
      | main    | git rebase origin/main             |
      |         | git push                           |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
    And it prints the error:
      """
      To abort, run "git-town sync --abort".
      To continue after you have resolved the conflicts, run "git-town sync --continue".
      To skip the sync of the 'feature' branch, run "git-town sync --skip".
      """
    And I am still on the "feature" branch
    And my uncommitted file is stashed
    And my repo has a merge in progress


  Scenario: aborting
    When I run `git-town sync --abort`
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                  |
      | feature | git merge --abort                                        |
      |         | git reset --hard <%= sha 'conflicting feature commit' %> |
      |         | git checkout main                                        |
      | main    | git checkout feature                                     |
      | feature | git stash pop                                            |
    And I am still on the "feature" branch
    And my workspace has the uncommitted file again
    And there is no merge in progress
    And my repository still has the following commits
      | BRANCH  | LOCATION         | MESSAGE                    | FILE NAME        | FILE CONTENT    |
      | main    | local and remote | conflicting main commit    | conflicting_file | main content    |
      | feature | local            | conflicting feature commit | conflicting_file | feature content |
      |         | remote           | feature commit             | feature_file     | feature content |


  Scenario: continuing without resolving the conflicts
    When I run `git-town sync --continue`
    Then Git Town runs no commands
    And it prints the error "You must resolve the conflicts before continuing"
    And I am still on the "feature" branch
    And my uncommitted file is stashed
    And my repo still has a merge in progress


  Scenario: continuing after resolving the conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `git-town sync --continue`
    Then Git Town runs the commands
      | BRANCH  | COMMAND              |
      | feature | git commit --no-edit |
      |         | git push             |
      |         | git stash pop        |
    And I am still on the "feature" branch
    And my workspace has the uncommitted file again
    And my repository still has the following commits
      | BRANCH  | LOCATION         | MESSAGE                                                    | FILE NAME        |
      | main    | local and remote | conflicting main commit                                    | conflicting_file |
      | feature | local and remote | conflicting feature commit                                 | conflicting_file |
      |         |                  | feature commit                                             | feature_file     |
      |         |                  | Merge remote-tracking branch 'origin/feature' into feature |                  |
      |         |                  | conflicting main commit                                    | conflicting_file |
      |         |                  | Merge branch 'main' into feature                           |                  |
    And my repository still has the following committed files
      | BRANCH  | NAME             | CONTENT          |
      | main    | conflicting_file | main content     |
      | feature | conflicting_file | resolved content |
      | feature | feature_file     | feature content  |


  Scenario: continuing after resolving the conflicts and comitting
    Given I resolve the conflict in "conflicting_file"
    When I run `git commit --no-edit; git-town sync --continue`
    Then Git Town runs the commands
      | BRANCH  | COMMAND       |
      | feature | git push      |
      |         | git stash pop |
    And I am still on the "feature" branch
    And my workspace has the uncommitted file again
    And my repository still has the following commits
      | BRANCH  | LOCATION         | MESSAGE                                                    | FILE NAME        |
      | main    | local and remote | conflicting main commit                                    | conflicting_file |
      | feature | local and remote | conflicting feature commit                                 | conflicting_file |
      |         |                  | feature commit                                             | feature_file     |
      |         |                  | Merge remote-tracking branch 'origin/feature' into feature |                  |
      |         |                  | conflicting main commit                                    | conflicting_file |
      |         |                  | Merge branch 'main' into feature                           |                  |
    And my repository still has the following committed files
      | BRANCH  | NAME             | CONTENT          |
      | main    | conflicting_file | main content     |
      | feature | conflicting_file | resolved content |
      | feature | feature_file     | feature content  |
