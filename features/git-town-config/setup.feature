Feature: Initial configuration

  As a user who hasn't configured Git Town yet
  I want to have a simple, dedicated setup command
  So that I can configure it safely before using any Git Town command


  Background:
    Given my repository has the branches "production" and "dev"
    And I haven't configured Git Town yet


  Scenario: succeeds on valid main branch and perennial branch names
    When I run `git-town config --setup` and enter:
      | main       |
      | production |
      | dev        |
      |            |
    Then the main branch is now configured as "main"
    And the perennial branches are now configured as "production" and "dev"


  Scenario: succeeds on valid main branch and perennial branch numbers
    When I run `git-town config --setup` and enter:
      | 2 |
      | 1 |
      | 3 |
      |   |
    Then the main branch is now configured as "main"
    And the perennial branches are now configured as "production" and "dev"


  Scenario: shows error and re-prompts on empty main branch
    When I run `git-town config --setup` and enter:
      |      |
      | main |
      |      |
    Then it prints "A main development branch is required to enable the features provided by Git Town"
    And the main branch is now configured as "main"
    And my repo is configured with no perennial branches


  Scenario: shows error and re-prompts on invalid main branch number
    When I run `git-town config --setup` and enter:
      | 4    |
      | main |
      |      |
    Then it prints "Invalid branch number"
    And the main branch is now configured as "main"
    And my repo is configured with no perennial branches


  Scenario: shows error and re-prompts on non-existent main branch
    When I run `git-town config --setup` and enter:
      | non-existent |
      | main         |
      |              |
    Then it prints "Branch 'non-existent' doesn't exist"
    And the main branch is now configured as "main"
    And my repo is configured with no perennial branches


  Scenario: shows error and re-prompts on main branch as perennial branch
    When I run `git-town config --setup` and enter:
      | main |
      | main |
      | dev  |
      |      |
    Then it prints "'main' is already set as the main branch"
    And the main branch is now configured as "main"
    And the perennial branches are now configured as "dev"


  Scenario: shows error and re-prompts on invalid perennial branch number
    When I run `git-town config --setup` and enter:
      | main |
      | 4    |
      | 3    |
      |      |
    Then it prints "Invalid branch number"
    And the main branch is now configured as "main"
    And the perennial branches are now configured as "production"


  Scenario: shows error and re-prompts on non-existent perennial branch
    When I run `git-town config --setup` and enter:
      | main         |
      | non-existent |
      | dev          |
      |              |
    Then it prints "Branch 'non-existent' doesn't exist"
    And the main branch is now configured as "main"
    And the perennial branches are now configured as "dev"
