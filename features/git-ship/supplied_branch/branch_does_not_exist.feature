Feature: git ship: errors when trying to ship a branch that doesn't exist

  As a developer trying to ship a branch that doesn't exist
  I should see an error telling me about this
  So that I can ship the correct branch and remain productive.


  Background:
    Given I am on the "main" branch


  Scenario: with open changes
    Given I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git ship non-existing-branch` while allowing errors
    Then it runs no Git commands
    And I get the error "There is no branch named 'non-existing-branch'"
    And I end up on the "main" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"


  Scenario: without open changes
    When I run `git ship non-existing-branch` while allowing errors
    Then it runs no Git commands
    And I get the error "There is no branch named 'non-existing-branch'"
    And I end up on the "main" branch
