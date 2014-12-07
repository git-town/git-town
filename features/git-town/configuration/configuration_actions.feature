Feature: listing the configuration

  Scenario: everything is configured
    Given I have configured the main branch name as "main"
    And my non-feature branches are configured as "qa, staging"
    When I run `git town --config`
    Then I see "Main branch: 'main'"
    And I see "Non-feature branches: 'qa, staging'"


  Scenario: the main branch is configured but the non-feature branches are not
    Given I have configured the main branch name as "main"
    And my non-feature branches are not configured
    When I run `git town --config`
    Then I see "Main branch: 'main'"
    And I see "Non-feature branches: ''"


  Scenario: the main branch is not configured but the non-feature branches are
    Given I don't have a main branch name configured
    And my non-feature branches are configured as "qa"
    When I run `git town --config`
    Then I see "Main branch: ''"
    And I see "Non-feature branches: 'qa'"


  Scenario: nothing is configured yet
    Given I haven't configured Git Town yet
    When I run `git town --config`
    Then I see "Main branch: ''"
    And I see "Non-feature branches: ''"


  Scenario: printing just the main branch when it's not yet configured
    Given I don't have a main branch name configured
    When I run `git town main-branch`
    Then I see "Main branch: ''"


  Scenario: printing just the main branch when it's configured
    Given I have configured the main branch name as "main"
    When I run `git town main-branch`
    Then I see "Main branch: 'main'"


  Scenario: printing just the non-feature branches when they're not yet configured
    Given my non-feature branches are not configured
    When I run `git town non-feature-branches`
    Then I see "Non-feature branches: ''"


  Scenario: printing just the non-feature branches when they're configured
    Given my non-feature branches are configured as "qa"
    When I run `git town non-feature-branches`
    Then I see "Non-feature branches: 'qa'"


  Scenario: setting the main branch
    Given I have configured the main branch name as "main-old"
    When I run `git town main-branch main-new`
    Then I see "main branch stored as 'main-new'"