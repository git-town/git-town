Feature: errors while shipping the supplied branch without open changes


  Scenario: feature branch that does not exist
    Given I am on the "feature" branch
    When I run `git ship other_feature -m 'feature done'` while allowing errors
    Then I get the error "There is no branch named 'other_feature'."
    And I end up on the "feature" branch


  Scenario: does not ship a feature branch that is not ahead of main
    Given I have feature branches named "feature" and "other_feature"
    And I am on the "other_feature" branch
    When I run `git ship feature -m 'feature done'` while allowing errors
    Then I get the error "The branch 'feature' has no commits to merge into 'main'."
    And I end up on the "other_feature" branch


  Scenario: does not ship the main branch
    Given I am on the "feature" branch
    When I run `git ship main -m 'feature done'` while allowing errors
    Then I get the error "The branch 'main' is not a feature branch. Only feature branches can be shipped."
    And I am still on the "feature" branch


  Scenario: does not ship a non-feature branch
    Given non-feature branch configuration "qa, production"
    And I am on the "feature" branch
    When I run `git ship production -m 'feature done'` while allowing errors
    Then I get the error "The branch 'production' is not a feature branch. Only feature branches can be shipped."
    And I am still on the "feature" branch
