Feature: Git Ship: errors when shipping the main branch without open changes

  Background:
    Given I am on the "feature" branch
    When I run `git ship main -m 'feature done'` while allowing errors


  Scenario: result
    Then I get the error "The branch 'main' is not a feature branch. Only feature branches can be shipped."
    And I am still on the "feature" branch
