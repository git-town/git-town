Feature: git-hack enforces unique branch names while starting a new feature

  As a developer
  I should not be able to add a new feature into an already existing feature branch
  So that feature branches remain focussed and code reviews effective


  Background:
    Given I have a feature branch named "existing_feature"
    And I am on the main branch
    When I run `git hack existing_feature` while allowing errors


  Scenario: result
    Then it runs no Git commands
    And I get the error "A branch named 'existing_feature' already exists"
    And I am still on the "main" branch
