Feature: Show clear error if trying to continue after executing a successful command

  As a developer accidentally trying to continue a command after it completed successfully
  I should see a friendly and descriptive message that the command can not be continued
  So that I don't experience any cryptic errors.


  Scenario: continuing after successful git-hack
    Given I had run `git town-hack new-feature`
    When I run `git town-hack --continue`
    Then I get the error "The last command finished successfully and cannot be continued"


  Scenario: continuing after successful git-kill
    Given I have a feature branch named "current-feature"
    And I am on the "current-feature" branch
    And I run `git town-kill`
    When I run `git town-kill --continue`
    Then I get the error "The last command finished successfully and cannot be continued"


  Scenario: continuing after successful git-prune-branches
    Given I have a feature branch named "deleted-feature"
    And the "deleted-feature" branch gets deleted on the remote
    And I am on the "main" branch
    And I run `git town-prune-branches`
    When I run `git town-prune-branches --continue`
    Then I get the error "The last command finished successfully and cannot be continued"


  Scenario: continuing after successful git-rename-branch
    Given I have a feature branch named "feature"
    And I am on the "main" branch
    And I run `git town-rename-branch feature renamed-feature`
    When I run `git town-rename-branch --continue`
    Then I get the error "The last command finished successfully and cannot be continued"


  Scenario: continuing after successful git-ship
    Given I have a feature branch named "current-feature"
    And the following commit exists in my repository
      | BRANCH          | FILE NAME    |
      | current-feature | feature_file |
    And I am on the "current-feature" branch
    And I run `git town-ship -m "feature done"`
    When I run `git town-ship --continue`
    Then I get the error "The last command finished successfully and cannot be continued"


  Scenario: continuing after successful git-sync
    Given I am on the "main" branch
    And the following commits exist in my repository
      | LOCATION | FILE NAME   |
      | local    | local_file  |
      | remote   | remote_file |
    And I run `git town-sync`
    When I run `git town-sync --continue`
    Then I get the error "The last command finished successfully and cannot be continued"
