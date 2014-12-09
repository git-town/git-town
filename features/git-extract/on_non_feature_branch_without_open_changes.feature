Feature: git-extract errors on a non-feature branch

  Background:
    Given non-feature branch configuration "qa, production"
    And I am on the "production" branch
    When I run `git extract refactor` while allowing errors


  Scenario: result
    Then it runs no Git commands
    And I get the error "The branch 'production' is not a feature branch. You must be on a feature branch in order to extract commits."
    And I am still on the "production" branch
