Feature: preserve the previous Git branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | current  | feature | main   | local, origin |
      | previous | feature | main   | local, origin |
    And the current branch is "current" and the previous branch is "previous"

  Scenario: current branch renamed
    When I run "git-town rename current new"
    Then the current branch is now "new"
    And the previous Git branch is still "previous"

  Scenario: previous branch renamed
    When I run "git-town rename previous new"
    Then the current branch is now "current"
    And the previous Git branch is now "new"
