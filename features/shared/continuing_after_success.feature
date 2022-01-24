Feature: Show clear error if trying to continue after executing a successful command

  Scenario: continuing after successful git-hack
    Given I run "git-town hack new-feature"
    When I run "git-town continue"
    Then it prints the error:
      """
      nothing to continue
      """

  Scenario: continuing after successful git-ship
    Given my repo has a feature branch named "current-feature"
    And the following commits exist in my repo
      | BRANCH          | FILE NAME    |
      | current-feature | feature_file |
    And I am on the "current-feature" branch
    And I run "git-town ship -m 'feature done'"
    When I run "git-town continue"
    Then it prints the error:
      """
      nothing to continue
      """

  Scenario: continuing after successful git-sync
    Given I am on the "main" branch
    And the following commits exist in my repo
      | LOCATION | FILE NAME   |
      | local    | local_file  |
      | remote   | remote_file |
    And I run "git-town sync"
    When I run "git-town continue"
    Then it prints the error:
      """
      nothing to continue
      """
