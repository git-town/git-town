Feature: display the perennial branch configuration

  Scenario: no perennial branches configured
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
