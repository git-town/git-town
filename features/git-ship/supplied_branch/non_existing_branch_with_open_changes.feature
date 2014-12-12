Feature: git ship: does not ship a non-existing branch (with open changes)

  As a developer mistyping the branch to be shipped
  I want to get a notification that the branch I provided doesn't exist
  So that I can correct my mistake, ship the correct branch, and remain productive.


  Background:
    Given I am on the "feature" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git ship other_feature -m 'feature done'` while allowing errors


  Scenario: result
    Then I get the error "There is no branch named 'other_feature'"
    And I end up on the "feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
