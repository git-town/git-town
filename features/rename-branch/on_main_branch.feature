Feature: git town-rename-branch: errors when renaming the main branch

  As a developer accidentally trying to rename the main branch
  I should see an error that this is not possible
  So that I know that only other branches can be renamed.


  Background:
    Given I am on the "main" branch


  Scenario: error when trying to rename
    When I run `git-town rename-branch main renamed-main`
    Then it runs no commands
    And I get the error "The main branch cannot be renamed."
    And I am still on the "main" branch


  Scenario: error when trying to force rename
    When I run `git-town rename-branch main renamed-main --force`
    Then it runs no commands
    And I get the error "The main branch cannot be renamed."
    And I am still on the "main" branch
