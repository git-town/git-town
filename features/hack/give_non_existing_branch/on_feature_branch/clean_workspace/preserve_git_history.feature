Feature: preserve the previous Git branch

  Scenario:
    Given a Git repo clone
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | current  | feature | main   | local     |
      | previous | feature | main   | local     |
    And the current branch is "current" and the previous branch is "previous"
    When I run "git-town hack new"
    Then the current branch is now "new"
    And the previous Git branch is now "current"
