Feature: display the non-feature branch configuration

  Scenario: non-feature branches are not configured
    Given my non-feature branches are not configured
    When I run `git town non-feature-branches`
    Then I see "[none]"


  Scenario: non-feature branches are configured
    Given my non-feature branch is "qa" and "production"
    When I run `git town non-feature-branches`
    Then I see
      """
      qa
      production
      """
