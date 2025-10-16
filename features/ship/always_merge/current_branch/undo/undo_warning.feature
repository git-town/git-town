Feature: git-town undo prints a warning message for a merge commit

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
    And Git setting "git-town.ship-strategy" is "always-merge"
    And the current branch is "feature"
    And I ran "git-town ship -m 'feature done'"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town prints something like:
      """
      Cannot undo commit ".*" because it is on a perennial branch
      """
