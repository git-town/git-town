@messyoutput
Feature: ask for missing configuration

  Scenario:
    Given a Git repo with origin
    And Git Town is not configured
    When I run "git-town delete" and enter into the dialog:
      | DIALOG             | KEYS  |
      | welcome            | enter |
      | aliases            | enter |
      | main branch        | enter |
      | perennial branches |       |
      | origin hostname    | enter |
      | forge type         | enter |
      | enter all          | enter |
      | config storage     | enter |
    Then Git Town prints the error:
      """
      you cannot delete the main branch
      """
    And the main branch is now "main"
