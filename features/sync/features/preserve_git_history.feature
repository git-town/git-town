Feature: preserve the previous Git branch

  Background:
    Given a Git repo clone
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | current  | feature | main   | local, origin |
      | previous | feature | main   | local, origin |
    And the current branch is "current" and the previous branch is "previous"
    When I run "git-town sync"

  Scenario: result
    Then the current branch is still "current"
    And the previous Git branch is still "previous"

  Scenario: undo
    When I run "git-town undo"
    Then the current branch is still "current"
    And the previous Git branch is now "previous"
