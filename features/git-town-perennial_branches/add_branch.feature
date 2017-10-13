Feature: add a branch to the perennial branches configuration

  As a user or tool configuring Git Town's perennial branches
  I want an easy way to add a branch to my set of perennial branches
  So that I can configure Git Town safely, and the tool does exactly what I want.


  Background:
    Given my repository has the branches "staging" and "qa"
    And Git Town's perennial branches are configured as "qa"


  Scenario: adding an existing branch
    When I run `git-town perennial-branches --add staging`
    Then Git Town prints no output
    And its perennial branches are now configured as "qa" and "staging"


  Scenario: adding a branch that is already a perennial branch
    When I run `git-town perennial-branches --add qa`
    Then Git Town prints the error "'qa' is already a perennial branch"


  Scenario: adding a branch that is already set as the main branch
    Given Git Town's main branch is configured as "staging"
    When I run `git-town perennial-branches --add staging`
    Then Git Town prints the error "'staging' is already set as the main branch"


  Scenario: adding a branch that does not exist
    When I run `git-town perennial-branches --add branch-does-not-exist`
    Then Git Town prints the error "no branch named 'branch-does-not-exist'"


  Scenario: not providing a branch name
    When I run `git-town perennial-branches --add`
    Then Git Town prints the error "Error: flag needs an argument: --add"
    And it prints the error:
      """
      Usage:
        git-town perennial-branches [flags]
      """
