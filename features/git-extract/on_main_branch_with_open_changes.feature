Feature: git-extract errors on the main branch with open changes

  Background:
    Given I am on the "main" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git extract refactor` while allowing errors


  Scenario: result
    Then it runs no Git commands
    And I get the error "The branch 'main' is not a feature branch. You must be on a feature branch in order to extract commits."
    And I am still on the "main" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
