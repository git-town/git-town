Feature: Git checkout history is preserved when renaming the current branch

  Scenario: rename-branch
    Given my repo has the feature branches "previous" and "current"
    And I am on the "current" branch with "previous" as the previous Git branch
    When I run "git-town rename-branch current current-new"
    Then I am now on the "current-new" branch
    And the previous Git branch is still "previous"
