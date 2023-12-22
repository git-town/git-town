Feature: reset the configuration

  Scenario: with configuration
    Given the main branch is "main"
    And the current branch is a feature branch "feature"
    And the perennial branches are "qa" and "staging"
    When I run "git-town config reset"
    Then Git Town is no longer configured

  Scenario: no configuration
    Given Git Town is not configured
    When I run "git-town config reset"
    Then Git Town is no longer configured
