Feature: non-feature branch configuration

  Scenario: printing just the non-feature branches when they're not yet configured
    Given my non-feature branches are not configured
    When I run `git town non-feature-branches`
    Then I see "Non-feature branches: [none]"


  Scenario: printing just the non-feature branches when they're configured
    Given my non-feature branches are configured as "qa"
    When I run `git town non-feature-branches`
    Then I see "Non-feature branches: qa"


  Scenario: adding a non-feature-branch that doesn't exist
    Given my non-feature branches are configured as "qa"
    When I run `git town non-feature-branches --add staging`
    Then I see "Added 'staging' as a non-feature branch"
    And the non-feature branches include "staging"


  Scenario: adding a non-feature-branch that already exists
    Given my non-feature branches are configured as "qa"
    When I run `git town non-feature-branches --add qa`
    Then I see "'qa' is already a non-feature branch"
    And the non-feature branches include "qa"


  Scenario: removing a non-feature-branch that doesn't exist
    Given my non-feature branches are configured as "staging, qa"
    When I run `git town non-feature-branches --remove stage`
    Then I see "'stage' is not a non-feature branch"
    And the non-feature branches include "staging"
    And the non-feature branches include "qa"


  Scenario: removing a non-feature-branch that exists
    Given my non-feature branches are configured as "staging, qa"
    When I run `git town non-feature-branches --remove staging`
    Then I see "Removed 'staging' from non-feature branches"
    And the non-feature branches don't include "staging"
    And the non-feature branches include "qa"
