Feature: git town-prepend: errors when trying to prepend something in front of the main branch

  As a developer accidentally trying to prepend someting in front of the main branch
  I should see an error that the main branch has no parents
  So that I know about my mistake and run "git hack" instead.

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
