Feature: show help screen when Git Town is configured

  As a user having forgotten how to use Git Town
  I want to see a helpful list of all commands
  So that I can refresh my memory quickly and move on to what I actually wanted to do.


  Background:
    Given the main branch is configured as "main"
    And the perennial branches are configured as "qa" and "staging"


  Scenario Outline:
    When I run `<COMMAND>`
    Then it prints
      """
      Usage:
        git-town [command]
      """

    Examples:
      | COMMAND       |
      | git-town      |
      | git-town help |
