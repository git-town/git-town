@messyoutput
Feature: ask for missing configuration information

  Scenario: unconfigured
    Given a Git repo with origin
    And Git Town is not configured
    And Git setting "git-town.ship-strategy" is "fast-forward"
    When I run "git-town ship" and enter into the dialog:
      | DIALOG             | KEYS  |
      | welcome            | enter |
      | aliases            | enter |
      | main branch        | enter |
      | perennial branches |       |
      | perennial regex    | enter |
      | origin hostname    | enter |
      | forge type         | enter |
      | enter all          | enter |
      | config storage     | enter |
    And the main branch is now "main"
    And Git Town prints the error:
      """
      cannot ship the main branch
      """
