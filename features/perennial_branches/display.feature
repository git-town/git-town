Feature: display the perennial branches configuration

  As a user or tool unsure about which branches are currently configured as the perennial branches
  I want to be able to see this information simply and directly
  So that I can use it without furter thinking or processing, and my Git Town workflows are effective.

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
