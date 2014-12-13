Feature: git ship: does not ship main branch

  As a developer accidentally trying to ship the main branch
  I want to be reminded that this is not possible
  So that I can ship the right things without having to read the manual.


  Background:
    Given I am on the main branch
    When I run `git ship -m 'something done'` while allowing errors


  Scenario: result
    Then I get the error "The branch 'main' is not a feature branch. Only feature branches can be shipped."
    And I am still on the "main" branch
    And there are no commits
    And there are no open changes
