Feature: Initial configuration

  As a user who hasn't configured Git Town yet
  I want to have a simple, dedicated setup command
  So that I can configure it safely before using any Git Town command


  Background:
    Given I haven't configured Git Town yet


  Scenario: user enters non-existent main branch
    When I run `git town config --setup` and enter main branch name "nonexistent"
    Then I get the error "There is no branch named 'nonexistent'"
    And Git Town is still not configured for this repository


  Scenario: user enters a valid main branch and non-existent perennial branch
    Given I have a branch named "master"
    When I run `git town config --setup` and enter main branch name "master" and perennial branch names "nonexistent"
    Then I get the error "There is no branch named 'nonexistent'"
    And the main branch name is now configured as "master"
    And my perennial branches are still not configured


  Scenario: user enters valid main branch and valid perennial branches
    Given I have branches named "dev" and "qa"
    When I run `git town config --setup` and enter main branch name "dev" and perennial branch names "qa"
    Then the main branch name is now configured as "dev"
    And my perennial branches are now configured as "qa"


  Scenario: user enters the existing main branch as feature branches
    Given I have branches named "dev" and "qa"
    When I run `git town config --setup` and enter main branch name "dev" and perennial branch names "dev"
    Then I get the error "'dev' is already set as the main branch"
    And the main branch name is now configured as "dev"
    And my perennial branches are still not configured
