Feature: git hack: enforces unique branch names while starting a new feature

  As a developer starting work on a new branch
  I should be told when the branch name is taken
  So that I don't mix features, code reviews are easy, and the team productivity remains high.


  Background:
    Given I have a feature branch named "feature"
    And I am on the main branch
    When I run `git hack feature` while allowing errors


  Scenario: result
    Then it runs no Git commands
    And I get the error "A branch named 'feature' already exists"
    And I am still on the "main" branch
