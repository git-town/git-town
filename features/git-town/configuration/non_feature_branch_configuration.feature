Feature: non-feature branch configuration

  Scenario: printing just the non-feature branches when they're not yet configured
    Given my non-feature branches are not configured
    When I run `git town non-feature-branches`
    Then I see "Non-feature branches: ''"


  Scenario: printing just the non-feature branches when they're configured
    Given my non-feature branches are configured as "qa"
    When I run `git town non-feature-branches`
    Then I see "Non-feature branches: 'qa'"


  Scenario: adding a non-feature-branch that doesn't exist
    Given my non-feature branches are configured as "qa"
    When I run `git town non-feature-branches --add staging`
    Then I see "Added 'staging'"
    And the non-feature branches include "staging"


  Scenario: adding a non-feature-branch that already exists
    Given my non-feature branches are configured as "qa"
    When I run `git town non-feature-branches --add qa`
    Then I see "'qa' is already a non-feature branch"
    And the non-feature branches include "qa"
