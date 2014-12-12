Feature: git ship: does not ship a non-existing branch (without open changes)

  As a developer accidentally trying to ship a non-feature branch
  I should be notified about my mistake
  So that I can ship the correct branch and remain productive.


  Background:
    Given I am on the "feature" branch
    When I run `git ship other_feature -m 'feature done'` while allowing errors


  Scenario: result
    Then I get the error "There is no branch named 'other_feature'"
    And I end up on the "feature" branch
