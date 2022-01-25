Feature: errors if renaming a feature branch that has unpulled changes

  Background:
    Given my repo has a feature branch named "current-feature"

  Scenario: unpulled remote commits
    And the following commits exist in my repo
      | BRANCH          | LOCATION | MESSAGE               |
      | current-feature | remote   | remote feature commit |
    And I am on the "current-feature" branch
    When I run "git-town rename-branch current-feature renamed-feature"
    Then it runs the commands
      | BRANCH          | COMMAND                  |
      | current-feature | git fetch --prune --tags |
    And it prints the error:
      """
      "current-feature" is not in sync with its tracking branch, please sync the branches before renaming
      """
    And I am still on the "current-feature" branch

  Scenario: unpushed local commits
    And the following commits exist in my repo
      | BRANCH          | LOCATION | MESSAGE              |
      | current-feature | local    | local feature commit |
    And I am on the "current-feature" branch
    When I run "git-town rename-branch current-feature renamed-feature"
    Then it runs the commands
      | BRANCH          | COMMAND                  |
      | current-feature | git fetch --prune --tags |
    And it prints the error:
      """
      "current-feature" is not in sync with its tracking branch, please sync the branches before renaming
      """
    And I am now on the "current-feature" branch
