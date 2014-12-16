Feature: git ship: does not ship non-feature branches

  As a developer accidentally trying to ship a non-feature branch
  I should be reminded that this is not possible
  So that I can ship the right things without having to read the manual, and can focus on real work.


  Background:
    Given non-feature branch configuration "qa, production"
    And I am on the "production" branch
    When I run `git ship -m 'feature done'` while allowing errors


  Scenario: result
    Then I get the error "The branch 'production' is not a feature branch. Only feature branches can be shipped."
    And I am still on the "production" branch
    And there are no commits
    And there are no open changes
