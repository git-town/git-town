Feature: Show explanation when trying to continue after a successful command

  Scenario: continuing after successful git-hack
    Given I run "git-town hack new-feature"
    When I run "git-town continue"
    Then it prints the error:
      """
      nothing to continue
      """

  Scenario: continuing after successful git-ship
    Given my repo has a feature branch named "current-feature"
    And I run "git-town ship current-feature -m done"
    When I run "git-town continue"
    Then it prints the error:
      """
      nothing to continue
      """

  Scenario: continuing after successful git-sync
    Given I run "git-town sync"
    When I run "git-town continue"
    Then it prints the error:
      """
      nothing to continue
      """
