Feature: making the current parked branch a feature branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE   | PARENT | LOCATIONS |
      | parked | (none) | main   | local     |
    And the current branch is "parked"
    And local Git setting "git-town-branch.parked.branchtype" is "parked"
    When I run "git-town hack"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints:
      """
      branch "parked" is now a feature branch
      """
    And branch "parked" now has type "feature"
    And local Git setting "git-town-branch.parked.branchtype" is now "feature"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And branch "parked" now has type "parked"
    And local Git setting "git-town-branch.parked.branchtype" is now "parked"
