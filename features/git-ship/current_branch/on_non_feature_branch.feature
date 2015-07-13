Feature: git ship: errors when trying to ship a perennial branch

  As a developer accidentally trying to ship a perennial branch
  I should see an error that this is not possible
  So that I know how to ship things correctly without having to read the manual.


  Background:
    Given I have branches named "qa" and "production"
    And my perennial branches are configured as "qa" and "production"
    And I am on the "production" branch
    When I run `git ship -m "feature done"`


  Scenario: result
    Then I get the error "The branch 'production' is not a feature branch. Only feature branches can be shipped."
    And I am still on the "production" branch
    And there are no commits
    And there are no open changes
