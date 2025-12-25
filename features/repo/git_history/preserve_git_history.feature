@skipWindows
Feature: preserve the previous Git branch

  Scenario:
    Given a Git repo with origin
    And the origin is "https://github.com/git-town/git-town.git"
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | current  | feature | main   | local, origin |
      | previous | feature | main   | local, origin |
    And the current branch is "current" and the previous branch is "previous"
    And tool "open" is installed
    When I run "git-town repo"
    Then the current branch is still "current"
    And the previous Git branch is still "previous"
