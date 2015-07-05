Feature: git ship: errors if there are open changes

  As a developer trying to ship a branch with uncommitted changes
  I should see an error that my branch is in an unfinished state
  So that my users don't experience half-baked features.


  Background:
    Given I have a feature branch named "feature"
    And I have an uncommitted file
    And I am on the "feature" branch
    When I run `git ship -m "feature done"`


  Scenario: result
    Then it runs no Git commands
    And I get the error "You cannot ship with uncommitted changes."
    And I am still on the "feature" branch
    And I still have my uncommitted file
    And there are no commits


  Scenario: undo
    When I run `git ship --undo`
    Then I get the error "Nothing to undo"
    And it runs no Git commands
    And I am still on the "feature" branch
    And I am left with my original commits
    And I still have my uncommitted file
