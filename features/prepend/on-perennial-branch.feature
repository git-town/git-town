Feature: git town-prepend: errors when trying to prepend something in front of the main branch


  Background:
    Given my repo has the perennial branches "qa" and "production"
    And I am on the "production" branch
    When I run "git-town prepend new-parent"


  Scenario: result
    Then it runs the commands
      | BRANCH     | COMMAND                  |
      | production | git fetch --prune --tags |
    And it prints the error:
      """
      the branch "production" is not a feature branch. Only feature branches can have parent branches
      """
    And I am still on the "production" branch
