@skipWindows
Feature: multi-platform support

  Scenario Outline: supported tool installed
    Given my repo's origin is "https://github.com/git-town/git-town.git"
    And the "<TOOL>" tool is installed
    When I run "git-town repo"
    Then "open" launches a new pull request with this url in my browser:
      """
      https://github.com/git-town/git-town
      """

    Examples:
      | TOOL     |
      | open     |
      | xdg-open |

  Scenario: no supported tool installed
    Given my repo's origin is "https://github.com/git-town/git-town.git"
    And no tool to open browsers is installed
    When I run "git-town repo"
    Then it prints:
      """
      Please open in a browser: https://github.com/git-town/git-town
      """
