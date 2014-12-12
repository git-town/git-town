Feature: git extract: errors on a non-feature branch

  As a developer accidentally running `git extract` on a non-feature branch
  I want to be reminded about running this command on a feature branch
  So that I can use Git Town correctly without having to memorize the syntax.


  Background:
    Given non-feature branch configuration "qa, production"
    And I am on the "production" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git extract refactor` while allowing errors


  Scenario: result
    Then it runs no Git commands
    And I get the error "The branch 'production' is not a feature branch. You must be on a feature branch in order to extract commits."
    And I am still on the "production" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
