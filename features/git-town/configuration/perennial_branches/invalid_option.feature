Feature: passing an invalid option to the perennial branch configuration

  As a user or tool configuring Git Town's perennial branches
  I want to know when I use an invalid option
  So that I can configure Git Town safely, and the tool does exactly what I want.


  Scenario: using invalid option
    When I run `gt perennial-branches --invalid-option`
    Then I get the error "Error: unknown flag: --invalid-option"
    And I get the error
      """
      Usage:
        gt perennial-branches [flags]
      """
