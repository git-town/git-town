Feature: display the non-feature branch configuration

  As a user or tool unsure about which branches are currently configured as the non-feature branches
  I want to be able to see this information simply and directly
  So that I can use it without furter thinking or processing, and my Git Town workflows are effective.


  Scenario: non-feature branches are not configured
    Given my non-feature branches are not configured
    When I run `git town non-feature-branches`
    Then I see "Non-feature branches: [none]"


  Scenario: non-feature branches are configured
    Given my non-feature branch is "qa"
    When I run `git town non-feature-branches`
    Then I see "Non-feature branches: qa"
