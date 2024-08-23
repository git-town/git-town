@skipWindows
Feature: print the URL when the browser crashes

  Background:
    Given a Git repo with origin
    And the branch
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the current branch is "feature"
    And the origin is "git@github.com:git-town/git-town"
    And a proposal for this branch exists at "https://github.com/git-town/git-town/pull/123"
    And tool "open" is broken
    When I run "git-town propose"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                                            |
      | feature | git fetch --prune --tags                           |
      | <none>  | looking for proposal online ... ok                 |
      |         | open https://github.com/git-town/git-town/pull/123 |
    And it prints:
      """
      Please open in a browser: https://github.com/git-town/git-town/pull/123
      """

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
