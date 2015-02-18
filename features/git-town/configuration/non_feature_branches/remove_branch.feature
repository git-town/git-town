Feature: remove a branch from the non-feature branches configuration

  As a user or tool configuring Git Town's non-feature branches
  I want an easy way to remove a branch from my set of non-feature branches
  So that I can configure Git Town safely, and the tool does exactly what I want.


  Background:
    Given my non-feature branches are configured as "staging" and "qa"


  Scenario: removing a branch that is a non feature branch
    When I run `git town non-feature-branches --remove staging`
    Then I see no output
    And my non-feature branches are now configured as "qa"


  Scenario: removing a branch that is not a non-feature branch
    When I run `git town non-feature-branches --remove feature`
    Then I get the error
      """
      error: 'feature' is not a non-feature branch
      """

  Scenario: not providing a branch name
    When I run `git town non-feature-branches --remove`
    Then I get the error
      """
      error: missing branch name
      usage: git town non-feature-branches (--add | --remove) <branchname>
      """
