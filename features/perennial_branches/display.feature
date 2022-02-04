Feature: display the perennial branches

  Scenario: unconfigured
    Given the perennial branches are not configured
    When I run "git-town perennial-branches"
    Then it prints:
      """
      [none]
      """

  Scenario: configured
    Given the perennial branches are "qa" and "production"
    When I run "git-town perennial-branches"
    Then it prints:
      """
      qa
      production
      """
