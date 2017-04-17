Feature: add a branch to the perennial branches configuration

  As a user or tool configuring Git Town's perennial branches
  I want an easy way to add a branch to my set of perennial branches
  So that I can configure Git Town safely, and the tool does exactly what I want.


  Background:
    Given I have branches named "staging" and "qa"
    And my perennial branches are configured as "qa"


  Scenario: adding an existing branch
    When I run `gt perennial-branches --add staging`
    Then I see no output
    And my repo is configured with perennial branches as "qa" and "staging"


  Scenario: adding a branch that is already a perennial branch
    When I run `gt perennial-branches --add qa`
    Then I get the error "'qa' is already a perennial branch"


  Scenario: adding a branch that is already set as the main branch
    Given I have configured the main branch name as "staging"
    When I run `gt perennial-branches --add staging`
    Then I get the error "'staging' is already set as the main branch"


  Scenario: adding a branch that does not exist
    When I run `gt perennial-branches --add branch-does-not-exist`
    Then I get the error "no branch named 'branch-does-not-exist'"


  Scenario: not providing a branch name
    When I run `gt perennial-branches --add`
    Then I get the error "Error: flag needs an argument: --add"
    And I get the error
      """
      Usage:
        gt perennial-branches [flags]
      """
