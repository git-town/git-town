@skipWindows
Feature: multi-platform support

  Background:
    Given a Git repo clone

  Scenario Outline: supported tool installed
    Given the origin is "https://github.com/git-town/git-town.git"
    And tool "<TOOL>" is installed
    When I run "git-town repo"
    Then "open" launches a new proposal with this url in my browser:
      """
      https://github.com/git-town/git-town
      """

    Examples:
      | TOOL     |
      | open     |
      | xdg-open |

  Scenario: no supported tool installed
    Given the origin is "https://github.com/git-town/git-town.git"
    And no tool to open browsers is installed
    When I run "git-town repo"
    Then it prints:
      """
      Please open in a browser: https://github.com/git-town/git-town
      """
