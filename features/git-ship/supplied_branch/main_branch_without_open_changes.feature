Feature: git ship: does not ship the main branch (without open changes)

  As a developer accidentally trying to ship the main branch
  I should be notified about my mistake
  So that I can ship the correct branch and remain productive.


  Background:
    Given I am on the "feature" branch
    When I run `git ship main -m 'feature done'` while allowing errors


  Scenario: result
    Then I get the error "The branch 'main' is not a feature branch. Only feature branches can be shipped."
    And I am still on the "feature" branch
