Feature: Initial configuration

  As a user who hasn't configured Git Town yet
  I want to have a simple, dedicated setup command
  So that I can configure it safely before using any Git Town command


  Background:
    Given I have branches named "production" and "dev"
    And I haven't configured Git Town yet


  Scenario: succeeds on valid main branch and perennial branch names
    When I run `git town config --setup` and enter "main", "production", "dev" and ""
    Then my repo is configured with the main branch as "main"
    And my repo is configured with perennial branches as "production" and "dev"


  Scenario: succeeds on valid main branch and perennial branch numbers
    When I run `git town config --setup` and enter "2", "1", "3" and ""
    Then my repo is configured with the main branch as "main"
    And my repo is configured with perennial branches as "production" and "dev"


  Scenario: shows error and re-prompts on empty main branch
    When I run `git town config --setup` and enter "", "main" and ""
    Then I see "no input received"
    And my repo is configured with the main branch as "main"
    And my repo is configured with no perennial branches


  Scenario: shows error and re-prompts on invalid main branch number
    When I run `git town config --setup` and enter "4", "main" and ""
    Then I see "invalid branch number"
    And my repo is configured with the main branch as "main"
    And my repo is configured with no perennial branches


  Scenario: shows error and re-prompts on non-existent main branch
    When I run `git town config --setup` and enter "non-existent", "main" and ""
    Then I see "branch 'non-existent' doesn't exist"
    And my repo is configured with the main branch as "main"
    And my repo is configured with no perennial branches


  Scenario: shows error and re-prompts on main branch as perennial branch
    When I run `git town config --setup` and enter "main", "main", "dev" and ""
    Then I see "'main' is already set as the main branch"
    And my repo is configured with the main branch as "main"
    And my repo is configured with perennial branches as "dev"


  Scenario: shows error and re-prompts on invalid perennial branch number
    When I run `git town config --setup` and enter "main", "4", "3" and ""
    Then I see "invalid branch number"
    And my repo is configured with the main branch as "main"
    And my repo is configured with perennial branches as "production"


  Scenario: shows error and re-prompts on non-existent perennial branch
    When I run `git town config --setup` and enter "main", "non-existent", "dev", and ""
    Then I see "branch 'non-existent' doesn't exist"
    And my repo is configured with the main branch as "main"
    And my repo is configured with perennial branches as "dev"
