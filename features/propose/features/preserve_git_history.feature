@skipWindows
Feature: preserve the previous Git branch

  Scenario:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | previous | feature | main   | local, origin |
      | current  | feature | main   | local, origin |
    And tool "open" is installed
    And the origin is "https://github.com/git-town/git-town.git"
    And a proposal for this branch does not exist
    And the current branch is "current" and the previous branch is "previous"
    When I run "git-town propose"
    Then the current branch is still "current"
    And the previous Git branch is still "previous"
