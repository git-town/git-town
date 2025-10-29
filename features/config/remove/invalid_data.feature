Feature: reset invalid configuration

  Scenario: sync-feature-strategy is invalid
    Given a Git repo with origin
    And the main branch is "main"
    And local Git setting "git-town.sync-feature-strategy" is "--help"
    When I run "git-town config remove"
    Then Git Town runs no commands
    And Git Town is no longer configured
    And local Git setting "git-town.sync-feature-strategy" now doesn't exist
