Feature: display GitHub token configuration

  Background:
    Given a Git repo with origin
    And Git setting "git-town.github-token" is "github-token"

  Scenario: masks the configured GitHub token
    When I run "git-town config"
    Then Git Town prints:
      """
        GitHub token: (configured)
      """
    And Git Town does not print "github-token"
