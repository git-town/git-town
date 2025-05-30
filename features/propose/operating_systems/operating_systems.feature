@skipWindows
Feature: support many browsers and operating systems

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the current branch is "feature"

  Scenario Outline:
    And the origin is "https://github.com/git-town/git-town.git"
    And tool "<TOOL>" is installed
    And a proposal for this branch does not exist
    When I run "git-town propose"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                              |
      | feature | git fetch --prune --tags                                             |
      |         | Looking for proposal online ... ok                                   |
      |         | <TOOL> https://github.com/git-town/git-town/compare/feature?expand=1 |

    Examples:
      | TOOL          |
      | open          |
      | xdg-open      |
      | cygstart      |
      | x-www-browser |
      | firefox       |
      | opera         |
      | mozilla       |
      | netscape      |

  Scenario: no supported tool installed
    And the origin is "https://github.com/git-town/git-town.git"
    And no tool to open browsers is installed
    When I run "git-town propose"
    Then Git Town prints:
      """
      Please open in a browser: https://github.com/git-town/git-town/compare/feature?expand=1
      """
