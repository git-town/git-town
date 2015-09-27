Feature: Converting old to new configuration

  As a user having existing configuration in the 0.6 format
  I want that Git Town automatically updates it to the 0.7 format
  So that I can update to the new Git Town version seamlessly.


  Scenario: existing non-feature branch configuration
    Given my non-feature branches are configured as "foo" and "bar"
    When I run `git town`
    Then I have no non-feature branch configuration
    And my repo is configured with perennial branches as "foo" and "bar"
