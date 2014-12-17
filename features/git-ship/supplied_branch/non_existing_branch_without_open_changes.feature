Feature: git ship: don't ship non-existing branches (without open changes)

  Background:
    Given I am on the "feature" branch
    When I run `git ship non-existing-branch -m 'feature done'` while allowing errors


  Scenario: result
    Then I get the error "There is no branch named 'non-existing-branch'"
    And I end up on the "feature" branch
