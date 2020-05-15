Feature: git town-diff-parent: diffing a given feature branch

  (see ../../current_branch/on_feature_branch/with_parent_branch.feature)


  Scenario: result
    Given my repository has a feature branch named "feature-1"
    And my repository has a feature branch named "feature-2" as a child of "feature-1"
    And I am on the "feature-2" branch
    When I run "git-town diff-parent feature-2"
    Then it runs the commands
      | BRANCH    | COMMAND                       |
      | feature-2 | git diff feature-1..feature-2 |
    And I am still on the "feature-2" branch
