Feature: Git Ship: errors when shipping the main branch

  Scenario: result
    Given I am on the main branch
    When I run `git ship -m 'feature done'` while allowing errors
    Then I get the error "The branch 'main' is not a feature branch. Only feature branches can be shipped."
    And I am still on the "main" branch
