Feature: modifying the non-feature branch configuration

  Scenario: adding a new non-feature-branch
    Given my non-feature branch is "qa"
    When I run `git town non-feature-branches --add staging`
    Then I see "Added 'staging' as a non-feature branch"
    And the non-feature branches include "staging"


  Scenario: adding a non-feature-branch that already exists
    Given my non-feature branch is "qa"
    When I run `git town non-feature-branches --add qa`
    Then I see "'qa' is already a non-feature branch"
    And the non-feature branches include "qa"


  Scenario: removing a non-feature-branch that doesn't exist
    Given my non-feature branches are "staging" and "qa"
    When I run `git town non-feature-branches --remove non-existing-branch`
    Then I see "'non-existing-branch' is not a non-feature branch"
    And the non-feature branches include "staging"
    And the non-feature branches include "qa"


  Scenario: removing a non-feature-branch that exists
    Given my non-feature branches are "staging" and "qa"
    When I run `git town non-feature-branches --remove staging`
    Then I see "Removed 'staging' from non-feature branches"
    And the non-feature branches don't include "staging"
    And the non-feature branches include "qa"
