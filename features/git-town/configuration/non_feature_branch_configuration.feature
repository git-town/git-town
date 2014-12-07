Feature: non-feature branch configuration

  Scenario: printing just the non-feature branches when they're not yet configured
    Given my non-feature branches are not configured
    When I run `git town non-feature-branches`
    Then I see "Non-feature branches: ''"


  Scenario: printing just the non-feature branches when they're configured
    Given my non-feature branches are configured as "qa"
    When I run `git town non-feature-branches`
    Then I see "Non-feature branches: 'qa'"


