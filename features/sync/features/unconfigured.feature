@messyoutput
Feature: ask for missing configuration information

  Scenario: run unconfigured
    Given a Git repo with origin
    And Git Town is not configured
    When I run "git-town sync" and enter into the dialog:
      | DIALOG             | KEYS  |
      | welcome            | enter |
      | aliases            | enter |
      | main branch        | enter |
      | perennial branches |       |
      | origin hostname    | enter |
      | forge type         | enter |
      | enter all          | enter |
      | config storage     | enter |
    And the main branch is now "main"
