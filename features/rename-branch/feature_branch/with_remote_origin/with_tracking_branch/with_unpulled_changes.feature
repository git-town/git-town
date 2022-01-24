Feature: git town-rename-branch: errors if renaming a feature branch that has unpulled changes

  Background:
    Given my repo has a feature branch named "current-feature"
    And the following commits exist in my repo
      | BRANCH          | LOCATION      | MESSAGE               |
      | main            | local, remote | main commit           |
      | current-feature | local, remote | feature commit        |
      |                 | remote        | remote feature commit |
    And I am on the "current-feature" branch
    And my workspace has an uncommitted file
    When I run "git-town rename-branch current-feature renamed-feature"

  Scenario: result
    Then it runs the commands
      | BRANCH          | COMMAND                  |
      | current-feature | git fetch --prune --tags |
    And it prints the error:
      """
      "current-feature" is not in sync with its tracking branch, please sync the branches before renaming
      """
    And I am now on the "current-feature" branch
    And my workspace still contains my uncommitted file
    And my repo is left with my original commits
