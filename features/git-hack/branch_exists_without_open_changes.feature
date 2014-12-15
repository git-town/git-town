Feature: git hack: enforces unique branch names while starting a new feature

  As a developer trying to start a new feature on an already existing branch
  I should see an error telling me that the branch name is taken
  So that my feature branches are focussed, code reviews easy, and the team productivity remains high.


  Background:
    Given I have a feature branch named "existing_feature"
    And I am on the main branch
    When I run `git hack existing_feature` while allowing errors


  Scenario: result
    Then it runs no Git commands
    And I get the error "A branch named 'existing_feature' already exists"
    And I am still on the "main" branch
