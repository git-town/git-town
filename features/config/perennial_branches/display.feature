Feature: display the perennial branches

  Scenario: unconfigured
    Given the perennial branches are not configured
    When I run "git-town config perennial-branches"
    Then it prints:
      """
      (not set)
      """

  Scenario: configured in local Git configuration
    Given the perennial branches are "qa" and "production"
    When I run "git-town config perennial-branches"
    Then it prints:
      """
      qa
      production
      """

  Scenario: configured in global Git configuration
    Given global Git Town setting "perennial-branches" is "qa production"
    When I run "git-town config perennial-branches"
    Then it prints:
      """
      qa
      production
      """

  Scenario: configured in config file
    Given local Git Town setting "perennial-branches" doesn't exist
    And the configuration file:
      """
      [branches]
        perennials = ["qa", "production"]
      """
    When I run "git-town config perennial-branches"
    Then it prints:
      """
      qa
      production
      """
