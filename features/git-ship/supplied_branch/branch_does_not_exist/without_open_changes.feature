Feature: Git Ship: branch does not exist without open changes

  Background:
    Given I am on the "feature" branch
    When I run `git ship other_feature -m 'feature done'` while allowing errors


  Scenario: result
    Then I get the error "There is no branch named 'other_feature'."
    And I end up on the "feature" branch
