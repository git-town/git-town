Feature: git-extract errors on the main branch

  As a developer accidentally running `git extract` on the main branch
  I want to be reminded about running this command on a feature branch
  So that I can use Git Town correctly without having to memorize the syntax.


  Background:
    Given I am on the main branch
    When I run `git extract refactor` while allowing errors


  Scenario: result
    Then I get the error "The branch 'main' is not a feature branch. You must be on a feature branch in order to extract commits."
    And I am still on the "main" branch
