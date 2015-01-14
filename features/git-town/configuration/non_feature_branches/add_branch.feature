Feature: add a branch to the non-feature branches configuration

  As a user or tool configuring Git Town's non-feature branches
  I want an easy way to add a branch to my set of non-feature branches
  So that I can configure Git Town safely, and the tool does exactly what I want.


  Background:
    Given I have branches named "staging" and "qa"
    And my non-feature branches are configured as "qa"


  Scenario: adding an existing branch
    When I run `git town non-feature-branches --add staging`
    Then I see no output
    And my non-feature branches are now configured as "qa" and "staging"


  Scenario: adding a branch that is already a non feature branch
    When I run `git town non-feature-branches --add qa` while allowing errors
    Then I see
      """
      error: 'qa' is already a non-feature branch
      """


  Scenario: adding a branch that does not exist
    When I run `git town non-feature-branches --add branch-does-not-exist` while allowing errors
    Then I see
      """
      error: no branch named 'branch-does-not-exist'
      """


  Scenario: not providing a branch name
    When I run `git town non-feature-branches --add` while allowing errors
    Then I see
      """
      error: missing branch name
      usage: git town non-feature-branches (--add | --remove) <branchname>
      """
