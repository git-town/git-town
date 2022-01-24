Feature: git town-diff-parent: diffing a given feature branch

  Scenario: result
    Given my repo has a feature branch named "feature-1"
    And my repo has a feature branch named "feature-2" as a child of "feature-1"
    And I am on the "feature-1" branch
    When I run "git-town diff-parent feature-2"
    Then it runs the commands
      | BRANCH    | COMMAND                       |
      | feature-1 | git diff feature-1..feature-2 |
    And I am still on the "feature-1" branch
