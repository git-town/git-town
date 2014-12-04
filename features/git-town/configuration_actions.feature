Feature: listing the configuration

  Scenario: everything is configured
    Given I have configured the main branch name as "main"
    And my non-feature branches are configured as "dev, qa, staging"
    When I run `git town --config`
    Then I see "Main branch: 'main'"
    And I see "Non-feature branches: 'dev, qa, staging'"


  Scenario: the main branch is configured but the non-feature branches are not
    Given I have configured the main branch name as "main"
    And my non-feature branches are not configured
    When I run `git town --config`
    Then I see "Main branch: 'main'"
    And I see "Non-feature branches: ''"


  Scenario: the main branch is not configured but the non-feature branches are
    Given I don't have a main branch name configured
    And my non-feature branches are configured as "dev"
    When I run `git town --config` and enter "master"
    Then I see "Main branch: 'master'"
    And I see "Non-feature branches: 'dev'"


  Scenario: nothing is configured yet
    Given I haven't configured Git Town yet
    When I run `git town --config` and enter "master" then enter "development"
    Then I see "Main branch: 'master'"
    And I see "Non-feature branches: 'development'"
