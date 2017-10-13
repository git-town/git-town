Feature: remove a branch from the perennial branches configuration

  As a user or tool configuring Git Town's perennial branches
  I want an easy way to remove a branch from my set of perennial branches
  So that I can configure Git Town safely, and the tool does exactly what I want.


  Background:
    Given Git Town's perennial branches are configured as "staging" and "qa"


  Scenario: removing a branch that is a perennial branch
    When I run `git-town perennial-branches --remove staging`
    Then Git Town prints no output
    And its perennial branches are now configured as "qa"


  Scenario: removing a branch that is not a perennial branch
    When I run `git-town perennial-branches --remove feature`
    Then Git Town prints the error "'feature' is not a perennial branch"


  Scenario: not providing a branch name
    When I run `git-town perennial-branches --remove`
    Then Git Town prints the error "Error: flag needs an argument: --remove"
    And it prints the error:
      """
      Usage:
        git-town perennial-branches [flags]
      """
