Feature: add perennial branches

  Background:
    Given the branches "staging" and "qa"

  @this
  Scenario: add a perennial branch to existing Git configuration
    Given the perennial branches are "qa"
    When I run "git-town config perennial-branches add staging"
    Then the perennial branches are now "qa" and "staging"
