@skipWindows
Feature: Create proposals for prototype branches

  Background:
    Given a Git repo clone
    And the branches
      | NAME      | TYPE      | PARENT | LOCATIONS     |
      | prototype | prototype | main   | local, origin |
    And the current branch is "prototype"
    And tool "open" is installed
    And the origin is "git@github.com:git-town/git-town.git"
    When I run "git-town propose"

  Scenario: result
    Then "open" launches a new proposal with this url in my browser:
      """
      https://github.com/git-town/git-town/compare/prototype?expand=1
      """
    And it prints:
      """
      branch "prototype" is no longer a prototype branch
      """
    And there are now no prototype branches
