Feature: git-extract errors on the main branch without open changes

  As a developer trying to extract commits from the main branch
  I want to be reminded that this command can only be run on feature branches
  So that I know how to do this correctly without having to memorize everything.


  Background:
    Given I am on the "main" branch
    When I run `git extract refactor` while allowing errors


  Scenario: result
    Then it runs no Git commands
    And I get the error "The branch 'main' is not a feature branch. You must be on a feature branch in order to extract commits."
    And I am still on the "main" branch
