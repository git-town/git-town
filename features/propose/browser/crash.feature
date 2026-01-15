@skipWindows
Feature: print the URL when the browser crashes

  Background:
    Given a Git repo with origin
    And the origin is "git@github.com:git-town/git-town"
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the proposals
      | ID | SOURCE BRANCH | TARGET BRANCH | URL                                           |
      |  1 | feature       | main          | https://github.com/git-town/git-town/pull/123 |
    And the current branch is "feature"
    And tool "open" is broken
    When I run "git-town propose"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                                        |
      | feature | git fetch --prune --tags                                                       |
      |         | Finding proposal from feature into main ... #1 (Proposal from feature to main) |
      |         | open https://github.com/git-town/git-town/pull/123                             |
    And Git Town prints:
      """
      Please open in a browser: https://github.com/git-town/git-town/pull/123
      """
    And the initial proposals exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the initial proposals exist now
