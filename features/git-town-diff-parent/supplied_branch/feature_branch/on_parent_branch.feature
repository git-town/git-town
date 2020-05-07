Feature: git town-diff-parent: diffing a given feature branch

  (see ../../current_branch/on_feature_branch/with_parent_branches.feature)


  Background:
    Given my repository has a feature branch named "feature-1"
    And my repository has a feature branch named "feature-2" as a child of "feature-1"
    And the following commits exist in my repository
      | BRANCH    | LOCATION      | MESSAGE          |
      | feature-1 | local, remote | feature 1 commit |
      | feature-2 | local, remote | feature 2 commit |
    And I am on the "feature-1" branch
    And my workspace has an uncommitted file
    When I run "git-town diff-parent feature-2"


  Scenario: result
    Then it runs the commands
      | BRANCH    | COMMAND                       |
      | feature-1 | git diff feature-1..feature-2 |
    And I am still on the "feature-1" branch
    And my workspace still contains my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES                   |
      | local      | main, feature-1, feature-2 |
      | remote     | main, feature-1, feature-2 |
    And my repository is left with my original commits
