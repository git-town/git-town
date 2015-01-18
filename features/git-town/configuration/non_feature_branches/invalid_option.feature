Feature: passing an invalid option to the non-feature branch configuration

  As a user or tool configuring Git Town's non-feature branches
  I want to know when I use an invalid option
  So that I can configure Git Town safely, and the tool does exactly what I want.


  Scenario: using invalid option
    When I run `git town non-feature-branches --invalid-option` it errors
    Then I see
      """
      error: unsupported option '--invalid-option'
      usage: git town non-feature-branches (--add | --remove) <branchname>
      """
