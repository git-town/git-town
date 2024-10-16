@messyoutput
Feature: ask for missing configuration

  Scenario:
    Given a Git repo with origin
    And Git Town is not configured
    When I run "git-town delete" and enter into the dialog:
      | DIALOG      | KEYS  |
      | main branch | enter |
    Then it prints the error:
      """
      you cannot delete the main branch
      """
    And the main branch is now "main"
