Feature: git ship: shipping the supplied feature branch with a tracking branch

  (see ../../current_branch/on_feature_branch/without_open_changes/with_tracking_branch.feature)


  Background:
    Given I have feature branches named "feature" and "other_feature"
    And the following commit exists in my repository
      | BRANCH  | LOCATION | MESSAGE        | FILE NAME    | FILE CONTENT    |
      | feature | remote   | feature commit | feature_file | feature content |
    And I am on the "other_feature" branch


  Scenario: with open changes
    Given I have an uncommitted file with name: "feature_file" and content: "conflicting content"
    When I run `git ship feature -m "feature done"`
    Then it runs the Git commands
      | BRANCH        | COMMAND                            |
      | other_feature | git stash -u                       |
      |               | git fetch --prune                  |
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
      |               | git checkout other_feature         |
      | other_feature | git stash pop                      |
    And I end up on the "other_feature" branch
    And my workspace still has an uncommitted file with name: "feature_file" and content: "conflicting content"
    And there is no "feature" branch
    And I have the following commits
      | BRANCH | LOCATION         | MESSAGE      | FILE NAME    |
      | main   | local and remote | feature done | feature_file |


  Scenario: without open changes
    When I run `git ship feature -m "feature done"`
    Then it runs the Git commands
      | BRANCH        | COMMAND                            |
      | other_feature | git fetch --prune                  |
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
      |               | git checkout other_feature         |
    And I end up on the "other_feature" branch
    And there is no "feature" branch
    And I have the following commits
      | BRANCH | LOCATION         | MESSAGE      | FILE NAME    |
      | main   | local and remote | feature done | feature_file |
