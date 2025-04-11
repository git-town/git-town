Feature: Automatically remove the "contribution_branches" setting

  Background:
    Given a local Git repo
    And local Git setting "git-town.contribution-branches" is "one two"
    When I run "git-town config"

  @this
  Scenario: result
    Then Git Town prints:
      """
      Dissolving deprecated local setting "git-town.contribution-branches"
      """
    And local Git setting "git-town.contribution-branches" now doesn't exist
    And local Git setting "git-town-branch.one.branchtype" is now "contribution"
