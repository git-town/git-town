Feature: git town-rename-branch: errors if renaming a perennial branch that has unpushed changes


  Background:
    Given my repo has the perennial branch "production"
    And the following commits exist in my repo
      | BRANCH     | LOCATION      | MESSAGE                  |
      | main       | local, remote | main commit              |
      | production | local, remote | production commit        |
      |            | local         | remote production commit |
    And I am on the "production" branch
    And my workspace has an uncommitted file
    When I run "git-town rename-branch --force production renamed-production"

  Scenario: result
    Then it runs the commands
      | BRANCH     | COMMAND                  |
      | production | git fetch --prune --tags |
    And it prints the error:
      """
      "production" is not in sync with its tracking branch, please sync the branches before renaming
      """
    And I am now on the "production" branch
    And my workspace still contains my uncommitted file
    And my repo is left with my original commits
