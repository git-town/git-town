Feature: git ship: don't ship non-existing branches (with open changes)

  As a developer mistyping the branch to be shipped
  I should be notified that the branch I provided doesn't exist
  So that I can correct my mistake, ship the correct branch, and remain productive.


  Background:
    Given I am on the "feature" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git ship non-existing-branch -m 'feature done'` while allowing errors


  Scenario: result
    Then I get the error "There is no branch named 'non-existing-branch'"
    And I end up on the "feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
