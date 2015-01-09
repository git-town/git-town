Feature: modifying the non-feature branch configuration

  Scenario: adding a new non-feature branch
    Given my non-feature branch is "qa"
    And I have a branch named "staging"
    When I run `git town non-feature-branches --add staging`
    Then I see no output
    And the non-feature branches include "staging"


  Scenario: adding a non-feature branch that already exists
    Given I have a branch named "qa"
    And my non-feature branch is "qa"
    When I run `git town non-feature-branches --add qa` while allowing errors
    Then I see
      """

        Error
        'qa' is already a non-feature branch

      """
    And the non-feature branches include "qa"


  Scenario: adding an invalid branch as a non-feature branch
    Given I have a branch named "qa"
    And my non-feature branch is "qa"
    When I run `git town non-feature-branches --add branch-does-not-exist` while allowing errors
    Then I see
      """

        Error
        There is no branch named 'branch-does-not-exist'

      """
    And the non-feature branches include "qa"


  Scenario: removing a non-feature branch that doesn't exist
    Given my non-feature branches are "staging" and "qa"
    When I run `git town non-feature-branches --remove non-existing-branch` while allowing errors
    Then I see
      """

        Error
        'non-existing-branch' is not a non-feature branch

      """
    And the non-feature branches include "staging"
    And the non-feature branches include "qa"


  Scenario: removing a non-feature branch that exists
    Given my non-feature branches are "staging" and "qa"
    When I run `git town non-feature-branches --remove staging`
    Then I see no output
    And the non-feature branches don't include "staging"
    And the non-feature branches include "qa"


  Scenario: using invalid parameter
    Given my non-feature branches are "staging" and "qa"
    When I run `git town non-feature-branches --invalid-parameter staging` while allowing errors
    Then I see
      """

        Error
        unknown option: --invalid-parameter

      """


  Scenario: missing branch name
    Given my non-feature branches are "staging" and "qa"
    When I run `git town non-feature-branches --add` while allowing errors
    Then I see
      """

        Error
        branch name required

      """
