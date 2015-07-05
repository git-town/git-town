Feature: git ship: errors when trying to ship a non-feature branch

  (see ../current_branch/on_non_feature_branch.feature)


  Background:
    Given I have branches named "qa" and "production"
    And my non-feature branches are configured as "qa" and "production"
    And I am on the "main" branch
    And I have an uncommitted file
    When I run `git ship production`


  Scenario: result
    Then it runs no Git commands
    And I get the error "The branch 'production' is not a feature branch. Only feature branches can be shipped."
    And I am still on the "main" branch
    And I still have my uncommitted file
