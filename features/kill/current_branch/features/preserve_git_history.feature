Feature: preserve the previous Git branch

  Background:
    Given a Git repo clone
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | current  | feature | main   | local, origin |
      | previous | feature | main   | local, origin |
    And the current branch is "current" and the previous branch is "previous"

  Scenario: previous branch remains
    When I run "git-town kill"
    Then the current branch is now "previous"
    And the previous Git branch is now ""

  Scenario: previous branch is gone
    When I run "git-town kill previous"
    Then the current branch is still "current"
    And the previous Git branch is now ""

  Scenario: current and previous branch remain
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | victim | feature | main   | local, origin |
    And the current branch is "current" and the previous branch is "previous"
    When I run "git-town kill victim"
    Then the current branch is still "current"
    And the previous Git branch is still "previous"
