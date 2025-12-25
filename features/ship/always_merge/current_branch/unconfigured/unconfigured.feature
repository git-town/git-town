@messyoutput
Feature: ask for missing configuration information

  Scenario: unconfigured
    Given a Git repo with origin
    And Git Town is not configured
    And Git setting "git-town.ship-strategy" is "always-merge"
    When I run "git-town ship" and enter into the dialog:
      | DIALOG             | KEYS  |
      | welcome            | enter |
      | aliases            | enter |
      | main branch        | enter |
      | perennial branches |       |
      | origin hostname    | enter |
      | forge type         | enter |
      | enter all          | enter |
      | config storage     | enter |
    And Git Town prints the error:
      """
      cannot ship the main branch
      """
    And the main branch is now "main"
