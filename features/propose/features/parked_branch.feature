@skipWindows
Feature: Create proposals for parked branches

  Background:
    Given a Git repo clone
    And the branches
      | NAME   | TYPE   | PARENT | LOCATIONS     |
      | parked | parked | main   | local, origin |
    And the current branch is "parked"
    And tool "open" is installed
    And the origin is "git@github.com:git-town/git-town.git"
    When I run "git-town propose"

  Scenario: result
    Then "open" launches a new proposal with this url in my browser:
      """
      https://github.com/git-town/git-town/compare/parked?expand=1
      """
