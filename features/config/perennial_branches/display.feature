Feature: display the perennial branches

  Scenario: unconfigured
    Given the perennial branches are not configured
    When I run "git-town config perennial-branches"
    Then it prints:
      """
      (not set)
      """

  Scenario: configured locally
    Given the perennial branches are "qa" and "production"
    When I run "git-town config perennial-branches"
    Then it prints:
      """
      qa
      production
      """
