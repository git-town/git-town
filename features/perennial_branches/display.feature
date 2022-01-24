Feature: display the perennial branches configuration

  Scenario: perennial branches are not configured
    Given the perennial branches are not configured
    When I run "git-town perennial-branches"
    Then it prints:
      """
      [none]
      """

  Scenario: perennial branches are configured
    Given the perennial branches are configured as "qa" and "production"
    When I run "git-town perennial-branches"
    Then it prints:
      """
      qa
      production
      """
