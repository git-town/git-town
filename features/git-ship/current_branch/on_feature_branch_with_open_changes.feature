Feature: git ship: don't ship unfinished features

  As a developer trying to ship a branch with uncommitted changes
  I should see an error that my branch is in an unfinished state
  So that my users don't experience half-baked features.


  Background:
    Given I have a feature branch named "feature"
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    And I am on the "feature" branch
    When I run `git ship -m "feature done"` it errors


  Scenario: result
    Then it runs no Git commands
    And I get the error "You cannot ship with uncommitted changes."
    And I am still on the "feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And there are no commits
