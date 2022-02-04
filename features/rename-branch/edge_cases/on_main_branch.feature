Feature: does not rename the main branch

  Background:
    Given I am on the "main" branch

  Scenario: try to rename
    When I run "git-town rename-branch main new"
    Then it runs no commands
    And it prints the error:
      """
      the main branch cannot be renamed
      """
    And I am still on the "main" branch

  Scenario: try to force rename
    When I run "git-town rename-branch main new --force"
    Then it runs no commands
    And it prints the error:
      """
      the main branch cannot be renamed
      """
    And I am still on the "main" branch
