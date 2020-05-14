Feature: git town-ship: shipping the current feature branch with a tracking branch

  (see ../../../current_branch/on_feature_branch/without_open_changes/with_tracking_branch.feature)


  Background:
    Given my repository has a feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION      | MESSAGE        | FILE NAME    | FILE CONTENT    |
      | feature | local, remote | feature commit | feature_file | feature content |
    And I am on the "feature" branch
    When I run "git-town ship feature -m 'feature done'"


  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune --tags           |
      |         | git checkout main                  |
      | main    | git rebase origin/main             |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
      |         | git checkout main                  |
      | main    | git merge --squash feature         |
      |         | git commit -m "feature done"       |
      |         | git push                           |
      |         | git push origin :feature           |
      |         | git branch -D feature              |
    And I end up on the "main" branch
    And the existing branches are
      | REPOSITORY | BRANCHES |
      | local      | main     |
      | remote     | main     |
    And my repository now has the following commits
      | BRANCH | LOCATION      | MESSAGE      | FILE NAME    |
      | main   | local, remote | feature done | feature_file |


  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                                       |
      | main    | git branch feature {{ sha 'feature commit' }} |
      |         | git push -u origin feature                    |
      |         | git revert {{ sha 'feature done' }}           |
      |         | git push                                      |
      |         | git checkout feature                          |
      | feature | git checkout main                             |
      | main    | git checkout feature                          |
    And I end up on the "feature" branch
    And my repository now has the following commits
      | BRANCH  | LOCATION      | MESSAGE               | FILE NAME    |
      | main    | local, remote | feature done          | feature_file |
      |         |               | Revert "feature done" | feature_file |
      | feature | local, remote | feature commit        | feature_file |
