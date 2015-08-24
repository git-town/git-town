Feature: git ship: errors when trying to ship a perennial branch

  (see ../current_branch/on_perennial_branch.feature)


  Background:
    Given I have branches named "qa" and "production"
    And my perennial branches are configured as "qa" and "production"
    And I am on the "main" branch
    And I have an uncommitted file
    When I run `git ship production`


  Scenario: result
    Then it runs no commands
    And I get the error "The branch 'production' is not a feature branch. Only feature branches can be shipped."
    And I am still on the "main" branch
    And I still have my uncommitted file
