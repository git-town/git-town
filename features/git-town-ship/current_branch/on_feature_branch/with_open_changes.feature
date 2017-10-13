Feature: git town-ship: errors if there are open changes

  As a developer trying to ship a branch with uncommitted changes
  I should see an error that my branch is in an unfinished state
  So that my users don't experience half-baked features.


  Background:
    Given my repository has a feature branch named "feature"
    And my workspace has an uncommitted file
    And I am on the "feature" branch
    When I run `git-town ship`


  Scenario: result
    Then Git Town runs no commands
    And it prints the error "You have uncommitted changes. Did you mean to commit them before shipping?"
    And I am still on the "feature" branch
    And my workspace still contains my uncommitted file
    And there are no commits


  Scenario: undo
    When I run `git-town ship --undo`
		Then Git Town runs no commands
    And it prints the error "Nothing to undo"
    And I am still on the "feature" branch
    And my workspace still contains my uncommitted file
