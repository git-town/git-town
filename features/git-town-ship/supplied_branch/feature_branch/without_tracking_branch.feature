Feature: git town-ship: shipping the supplied feature branch without a tracking branch

  (see ../../current_branch/on_feature_branch/without_open_changes/with_tracking_branch.feature)


  Background:
    Given my repository has the feature branches "feature" and "other-feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE        | FILE NAME    | FILE CONTENT    |
      | feature | local    | feature commit | feature_file | feature content |
    And I am on the "other-feature" branch
    And my workspace has an uncommitted file with name: "feature_file" and content: "conflicting content"
    When I run "git-town ship feature -m 'feature done'"


  Scenario: result
    Then it runs the commands
      | BRANCH        | COMMAND                            |
      | other-feature | git fetch --prune --tags           |
      |               | git add -A                         |
      |               | git stash                          |
      |               | git checkout main                  |
      | main          | git rebase origin/main             |
      |               | git checkout feature               |
      | feature       | git merge --no-edit origin/feature |
      |               | git merge --no-edit main           |
      |               | git checkout main                  |
      | main          | git merge --squash feature         |
      |               | git commit -m "feature done"       |
      |               | git push                           |
      |               | git push origin :feature           |
      |               | git branch -D feature              |
      |               | git checkout other-feature         |
      | other-feature | git stash pop                      |
    And I end up on the "other-feature" branch
    And my workspace still contains my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES            |
      | local      | main, other-feature |
      | remote     | main, other-feature |
    And my repository now has the following commits
      | BRANCH | LOCATION      | MESSAGE      | FILE NAME    |
      | main   | local, remote | feature done | feature_file |
