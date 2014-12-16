Feature: git ship: don't ship non-feature branches (with open changes)

  Background:
    Given non-feature branch configuration "qa, production"
    And I am on the "feature" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git ship production -m 'feature done'` while allowing errors


  Scenario: result
    Then I get the error "The branch 'production' is not a feature branch. Only feature branches can be shipped."
    And I am still on the "feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"

