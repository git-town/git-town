Feature: Show clear error if trying to continue after executing a successful command

  As a developer accidentally trying to continue a command after it completed successfully
  I should see a friendly and descriptive message that the command can not be continued
  So that I don't experience any cryptic errors.


  Scenario: continuing after successful git-hack
    Given I run `git-town hack new-feature`
    When I run `git-town continue`
    Then it prints the error "Nothing to continue"


  Scenario: continuing after successful git-ship
    Given my repository has a feature branch named "current-feature"
    And the following commit exists in my repository
      | BRANCH          | FILE NAME    |
      | current-feature | feature_file |
    And I am on the "current-feature" branch
    And I run `git-town ship -m "feature done"`
    When I run `git-town continue`
    Then it prints the error "Nothing to continue"


  Scenario: continuing after successful git-sync
    Given I am on the "main" branch
    And the following commits exist in my repository
      | LOCATION | FILE NAME   |
      | local    | local_file  |
      | remote   | remote_file |
    And I run `git-town sync`
    When I run `git-town continue`
    Then it prints the error "Nothing to continue"
