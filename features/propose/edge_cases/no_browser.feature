@skipWindows
Feature: print the URL when no browser installed

  Background:
    Given a Git repo with origin
    And the branch
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the origin is "git@github.com:git-town/git-town"
    And a proposal for this branch exists at "https://github.com/git-town/git-town/pull/123"
    And no tool to open browsers is installed
    And the current branch is "feature"
    When I run "git-town propose"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune --tags           |
      | <none>  | looking for proposal online ... ok |
    And it prints:
      """
      Please open in a browser: https://github.com/git-town/git-town/pull/123
      """

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
