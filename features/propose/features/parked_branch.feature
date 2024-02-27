@skipWindows
Feature: Create proposals for parked branches

  Background:
    Given the current branch is a parked branch "parked"
    And tool "open" is installed
    And the origin is "git@github.com:git-town/git-town.git"
    When I run "git-town propose"

  Scenario: result
    Then "open" launches a new proposal with this url in my browser:
      """
      https://github.com/git-town/git-town/compare/parked?expand=1
      """
