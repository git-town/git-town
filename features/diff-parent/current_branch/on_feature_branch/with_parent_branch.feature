Feature: git town-parent-diff: diffing the current feature branch

    As a developer makeing changes to a feature branched created with Git Town
    I want that Git Town tells me which changes this feature branch contains
    So that I know what work is left to do
    And which work will be merged into the parent branch

  Scenario: result
    Given my repo has a feature branch named "feature-1"
    And my repo has a feature branch named "feature-2" as a child of "feature-1"
    And I am on the "feature-2" branch
    When I run "git-town diff-parent"
    Then it runs the commands
      | BRANCH    | COMMAND                       |
      | feature-2 | git diff feature-1..feature-2 |
    And I am still on the "feature-2" branch
