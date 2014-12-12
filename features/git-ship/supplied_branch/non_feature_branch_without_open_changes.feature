Feature: git ship: does not ship a non-feature branch (without open changes)

  As a developer accidentally trying to ship a non-feature branch
  I should be notified about my mistake
  So that I can ship the correct branch and remain productive.


  Background:
    Given non-feature branch configuration "qa, production"
    And I am on the "feature" branch
    When I run `git ship production -m 'feature done'` while allowing errors


  Scenario: result
    Then I get the error "The branch 'production' is not a feature branch. Only feature branches can be shipped."
    And I am still on the "feature" branch
