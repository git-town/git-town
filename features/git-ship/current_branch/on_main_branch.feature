Feature: git ship: does not ship main branch

  Background:
    Given I am on the main branch
    When I run `git ship -m 'something done'` while allowing errors


  Scenario: result
    Then I get the error "The branch 'main' is not a feature branch. Only feature branches can be shipped."
    And I am still on the "main" branch
    And there are no commits
    And there are no open changes
