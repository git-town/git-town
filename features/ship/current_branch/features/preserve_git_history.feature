Feature: preserve the previous Git branch

  Background:
    Given a Git repo clone
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | previous | feature | main   | local, origin |
      | current  | feature | main   | local, origin |

  Scenario: current branch gone
    And the commits
      | BRANCH  | LOCATION |
      | current | local    |
    And the current branch is "current" and the previous branch is "previous"
    When I run "git-town ship -m 'feature done'"
    Then the current branch is now "main"
    And the previous Git branch is now "previous"

  Scenario: previous branch gone
    Given the commits
      | BRANCH   | LOCATION |
      | previous | local    |
    And the current branch is "current" and the previous branch is "previous"
    When I run "git-town ship previous -m 'feature done'"
    Then the current branch is still "current"
    And the previous Git branch is now "main"

  Scenario: both branches exist
    And the branch
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION |
      | feature | local    |
    And the current branch is "current" and the previous branch is "previous"
    When I run "git-town ship feature -m "feature done""
    Then the current branch is still "current"
    And the previous Git branch is still "previous"
