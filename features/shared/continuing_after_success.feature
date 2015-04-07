Feature: Show clear error if trying to continue after executing a successful command

  As a developer accidentally trying to continue a command after it completed successfully
  I should see a clear error that the command cannot be continued
  So that I don't experience any cryptic errors.


  Background:
    Given I have a feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION         |
      | feature | local and remote |
    And I run `git ship -m 'feature done'`


  Scenario: continuing after successful ship
    Then I run `git ship --continue`
    And I get the error "The last command finished successfully and cannot be continued"
