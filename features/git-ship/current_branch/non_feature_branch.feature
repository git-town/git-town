Feature: Git Ship: errors when shipping a non-feature branch

  Scenario: result
    Given non-feature branch configuration "qa, production"
    And I am on the "production" branch
    When I run `git ship -m 'feature done'` while allowing errors
    Then I get the error "The branch 'production' is not a feature branch. Only feature branches can be shipped."
    And I am still on the "production" branch
