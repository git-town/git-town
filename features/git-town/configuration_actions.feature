Feature: git town: listing the configuration

  As a user unsure about how Git Town is configured
  I want to be able to see the complete Git Town configuration with one command
  So that I can configure Git Town efficiently, and have more time for actual work.


  Scenario: everything is configured
    Given I have configured the main branch name as "main"
    And my non-feature branches are configured as "qa, staging"
    When I run `git town --config`
    Then I see "Main branch: 'main'"
    And I see "Non-feature branches: 'qa, staging'"


  Scenario: the main branch is configured but the non-feature branches are not
    Given I have configured the main branch name as "main"
    And my non-feature branches are not configured
    When I run `git town --config` and enter "production"
    Then I see "Main branch: 'main'"
    And I see "Non-feature branches: 'production'"


  Scenario: the main branch is not configured but the non-feature branches are
    Given I don't have a main branch name configured
    And my non-feature branches are configured as "qa"
    When I run `git town --config` and enter "master"
    Then I see "Main branch: 'master'"
    And I see "Non-feature branches: 'qa'"


  Scenario: nothing is configured yet
    Given I haven't configured Git Town yet
    When I run `git town --config` and enter "master" and "staging"
    Then I see "Main branch: 'master'"
    And I see "Non-feature branches: 'staging'"
